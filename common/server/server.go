package server

import (
	"net"
	"sync"
	"context"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"time"
	"os"
)

type options struct {
	codec      Codec
	hanlder    Handler
	onConnect  onConnectFunc
	onMessage  onMessageFunc
	onClose    onCloseFunc
	onError    onErrorFunc
	workerSize int // numbers of worker go-routines
	bufferSize int // size of buffered channel
}

// ServerOption sets server options.
type ServerOption func(*options)

// ReconnectOption returns a ServerOption that will make ClientConn reconnectable.
//func ReconnectOption() ServerOption {
//	return func(o *options) {
//		o.reconnect = true
//	}
//}

// 设置handler func
func CustomHandlerOption(handler Handler) ServerOption {
	return func(o *options) {
		o.hanlder = handler
	}
}

// CustomCodecOption returns a ServerOption that will apply a custom Codec.
func OnCustomCodecOption(codec Codec) ServerOption {
	return func(o *options) {
		o.codec = codec
	}
}

// TLSCredsOption returns a ServerOption that will set TLS credentials for server
// connections.
//func TLSCredsOption(config *tls.Config) ServerOption {
//	return func(o *options) {
//		o.tlsCfg = config
//	}
//}

// WorkerSizeOption returns a ServerOption that will set the number of go-routines
// in WorkerPool.
func WorkerSizeOption(workerSz int) ServerOption {
	return func(o *options) {
		o.workerSize = workerSz
	}
}

// BufferSizeOption returns a ServerOption that is the size of buffered channel,
// for example an indicator of BufferSize256 means a size of 256.
func BufferSizeOption(indicator int) ServerOption {
	return func(o *options) {
		o.bufferSize = indicator
	}
}

// OnConnectOption returns a ServerOption that will set callback to call when new
// client connected.
func OnConnectOption(cb func(WriteCloser) bool) ServerOption {
	return func(o *options) {
		o.onConnect = cb
	}
}

// OnMessageOption returns a ServerOption that will set callback to call when new
// message arrived.
func OnMessageOption(cb func(interface{}, WriteCloser)) ServerOption {
	return func(o *options) {
		o.onMessage = cb
	}
}

// OnCloseOption returns a ServerOption that will set callback to call when client
// closed.
func OnCloseOption(cb func(WriteCloser)) ServerOption {
	return func(o *options) {
		o.onClose = cb
	}
}

// OnErrorOption returns a ServerOption that will set callback to call when error
// occurs.
func OnErrorOption(cb func(WriteCloser)) ServerOption {
	return func(o *options) {
		o.onError = cb
	}
}

type server struct{
	opts          options
	ctx           context.Context
	cancel        context.CancelFunc
	conns         *sync.Map
	timing        *TimingWheel
	wg            *sync.WaitGroup
	mu            sync.Mutex
	netIdentifier *AtomicInt64
}

// ConnsSize returns connections size.
func (s *server) ConnsSize() int {
	var sz int
	s.conns.Range(func(k, v interface{}) bool {
		sz++
		return true
	})
	return sz
}

func newSuperServer(opts options) server {
	return server{
		opts:          opts,
		conns:         &sync.Map{},
		wg:            &sync.WaitGroup{},
		netIdentifier: NewAtomicInt64(0),
	}
}


// TCP Server
type TCPServer struct {
	server
	lis           map[net.Listener]bool			// TCP listens
}

// 创建tcpServer
func NewTCPServer(opt ...ServerOption) (*TCPServer, error) {
	var opts options
	for _, o := range opt {
		o(&opts)
	}

	if opts.codec == nil {
		return nil, ErrServerNeedCodec
	}
	if opts.workerSize <= 0 {
		opts.workerSize = defaultWorkersNum
	}
	if opts.bufferSize <= 0 {
		opts.bufferSize = BufferSize256
	}

	// initiates go-routine pool instance
	//globalWorkerPool = newWorkerPool(opts.workerSize)
	s := &TCPServer{
		server:newSuperServer(opts),
		lis:           make(map[net.Listener]bool),

	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.timing = NewTimingWheel(s.ctx)
	return s, nil
}

//func (s *Server) deleteUdpConn(l net)

// start tcp server
func (s *TCPServer) Start(l net.Listener) error {
	s.mu.Lock()
	if s.lis == nil {
		s.mu.Unlock()
		l.Close()
		return ErrServerClosed
	}
	s.lis[l] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		if s.lis != nil && s.lis[l] {
			l.Close()
			delete(s.lis, l)
		}
		s.mu.Unlock()
	}()

	xlogger.Infof("tcp server start, net %s addr %s\n", l.Addr().Network(), l.Addr().String())

	// TODO 处理timeout

	var tempDelay time.Duration
	for {
		tempDelay = 0
		rawConn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay >= max {
					tempDelay = max
				}
				xlogger.Errorf("accept error %v, retrying in %d\n", err, tempDelay)
				select {
				case <-time.After(tempDelay):
				case <-s.ctx.Done():
				}
				continue
			}
			return err
		}

		// how many connections do we have ?
		sz := s.ConnsSize()
		if sz >= MaxConnections {
			xlogger.Warnf("max connections size %d, refuse\n", sz)
			rawConn.Close()
			continue
		}

		// TODO TLS

		netId := s.netIdentifier.GetAndIncrement()
		// TODO newServerConn
		sc := NewTcpServerConn(netId, s, rawConn)

		// TODO sched?

		s.conns.Store(netId, sc)

		// TODO 分析??

		s.wg.Add(1)
		go func() {
			sc.Start()
		}()
		xlogger.Infof("accepted client %s, id %d, total %d\n", sc.GetName(), netId, s.ConnsSize())
		// TODO 打印连接信息?
		//s.conns.Range(func(k,v interface{}) bool{
		//	i := k.(int64)
		//	c := v.(*ServerConn)
		//	holmes.Infof("client(%d) %s", i, c.Name())
		//	return true
		//})

	}

	return nil
}


// Stop gracefully closes the server, it blocked until all connections
// are closed and all go-routines are exited.
func (s *TCPServer) Stop() {
	// immediately stop accepting new clients
	s.mu.Lock()
	listeners := s.lis
	s.lis = nil
	s.mu.Unlock()

	for l := range listeners {
		l.Close()
		xlogger.Infof("stop accepting at address %s\n", l.Addr().String())
	}

	// close all connections
	conns := map[int64]*ServerConn{}

	s.conns.Range(func(k, v interface{}) bool {
		i := k.(int64)
		c := v.(*ServerConn)
		conns[i] = c
		return true
	})
	// let GC do the cleanings
	s.conns = nil

	for _, c := range conns {
		c.rawConn.Close()
		xlogger.Infof("close client %s\n", c.GetName())
	}

	s.mu.Lock()
	s.cancel()
	s.mu.Unlock()

	s.wg.Wait()

	xlogger.Info("server stopped gracefully, bye.")
	os.Exit(0)
}


// udp Server
type UDPServer struct {
	server
}


// start udp server

func NewUDPServer(opt ...ServerOption) (*UDPServer, error) {
	var opts options
	for _, o := range opt {
		o(&opts)
	}

	if opts.codec == nil {
		return nil, ErrServerNeedCodec
	}
	if opts.workerSize <= 0 {
		opts.workerSize = defaultWorkersNum
	}
	if opts.bufferSize <= 0 {
		opts.bufferSize = BufferSize256
	}

	// initiates go-routine pool instance
	//globalWorkerPool = newWorkerPool(opts.workerSize)
	s := &UDPServer{
		server:newSuperServer(opts),

	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.timing = NewTimingWheel(s.ctx)
	return s, nil
}

func (s *UDPServer) Start(l *net.UDPConn) error {
	xlogger.Infof("udp server start, net %s\n", l.LocalAddr())
	// TODO 处理timeout
	// TODO TLS
	netId := s.netIdentifier.GetAndIncrement()
	// TODO newServerConn
	sc := NewUdpServerConn(netId, s, l)

	// TODO sched?

	s.conns.Store(netId, sc)

	// TODO 分析??

	s.wg.Add(1)
	go func() {
		sc.Start()
	}()
	xlogger.Infof("start upd lis %s, id %d, total %d\n", sc.GetName(), netId, s.ConnsSize())
	// TODO 打印连接信息?

	return nil
}

func (s *UDPServer) Stop() {
	// close all connections
	conns := map[int64]*UdpServerConn{}

	s.conns.Range(func(k, v interface{}) bool {
		i := k.(int64)
		c := v.(*UdpServerConn)
		conns[i] = c
		return true
	})
	for _, c := range conns {
		c.rawConn.Close()
		xlogger.Infof("close client %s\n", c.GetName())
	}
}