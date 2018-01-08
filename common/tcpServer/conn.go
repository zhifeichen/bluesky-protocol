package tcpServer

import (
	"bufio"
	"context"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"net"
	"sync"
)

// WriteCloser is the interface that groups Write and Close methods.
type WriteCloser interface {
	Write(interface{}) error
	Close()
}

type ServerConn struct {
	netId   int64
	belong  *Server
	rawConn net.Conn
	name    string

	once      *sync.Once
	wg        *sync.WaitGroup
	sendCh    chan []byte
	handlerCh chan interface{}
	//timerCh   chan *OnTimeOut

	mu     sync.Mutex // guards following
	ctx    context.Context
	cancel context.CancelFunc
}

func NewServerConn(id int64, s *Server, c net.Conn) *ServerConn {
	sc := &ServerConn{
		netId:   id,
		belong:  s,
		rawConn: c,
		once:    &sync.Once{},
		wg:      &sync.WaitGroup{},
		sendCh:    make(chan []byte, s.opts.bufferSize),
		handlerCh: make(chan interface{}, s.opts.bufferSize),
	}

	sc.ctx, sc.cancel = context.WithCancel(context.WithValue(s.ctx, serverCtx, s))

	sc.name = c.RemoteAddr().String()

	return sc
}

func (sc *ServerConn) SetName(name string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.name = name
}

func (sc *ServerConn) GetName() string {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.name
}

// Start starts the server connection, creating go-routines for reading,
// writing and handlng.
func (sc *ServerConn) Start() {
	xlogger.Infof("conn start, <%v -> %v>\n", sc.rawConn.LocalAddr(), sc.rawConn.RemoteAddr())
	onConnect := sc.belong.opts.onConnect
	if onConnect != nil {
		onConnect(sc)
	}

	loopers := []func(WriteCloser, *sync.WaitGroup){readLoop, writeLoop, handleLoop}

	for _, l := range loopers {
		sc.wg.Add(1)
		go l(sc, sc.wg)
	}

}

// RemoteAddr returns the remote address of server connection.
func (sc *ServerConn) RemoteAddr() net.Addr {
	return sc.rawConn.RemoteAddr()
}

// Close gracefully closes the server connection. It blocked until all sub
// go-routines are completed and returned.
func (sc *ServerConn) Close() {
	sc.once.Do(func() {
		xlogger.Infof("conn close gracefully, <%v -> %v>\n", sc.rawConn.LocalAddr(), sc.rawConn.RemoteAddr())

		// callback on close
		onClose := sc.belong.opts.onClose
		if onClose != nil {
			onClose(sc)
		}

		// remove connect from server
		sc.belong.conns.Delete(sc.netId)

		// TODO 分析?
		//addTotalConn(-1)

		// close net.Conn, any blocked read or write operation will be unblocked and
		// return errors.
		if tc, ok := sc.rawConn.(*net.TCPConn); ok {
			// avoid time-wait state
			// TCP将丢弃保留在套接口发送缓冲区中的任何数据并发送一个RST给对方，而不是通常的四分组终止序列，这避免了TIME_WAIT状态
			tc.SetLinger(0)
		}
		sc.rawConn.Close()

		// cancel readLoop, writeLoop and handleLoop go-routines.
		sc.mu.Lock()
		sc.cancel()

		sc.mu.Unlock()
		// clean up pending timers
		//for _, id := range pending {
		//	sc.CancelTimer(id)
		//}

		// wait until all go-routines exited.
		sc.wg.Wait()

		// close all channels and block until all go-routines exited.
		close(sc.sendCh)
		close(sc.handlerCh)
		//close(sc.timerCh)
		sc.belong.wg.Done()
	})
}

// Write writes a message to the client.
func (sc *ServerConn) Write(msg interface{}) error {
	return asyncWrite(sc, msg)

}

func asyncWrite(c WriteCloser, msg interface{}) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = ErrServerClosed
		}
	}()

	var (
		pkt    []byte
		sendCh chan []byte
	)
	sc := c.(*ServerConn)
	pkt, err = sc.belong.opts.codec.Encode(msg)
	sendCh = sc.sendCh

	if err != nil {
		xlogger.Errorf("asyncWrite error %v\n", err)
		return
	}

	select {
	case sendCh <- pkt:
		err = nil
	default:
		err = ErrWouldBlock
	}
	return
}

/**
1. scan and read data
2. do codec.decode([]byte) to msg
3. put msg to handle chan
4. wait close state
*/
func readLoop(c WriteCloser, wg *sync.WaitGroup) {

	var (
		rawConn   net.Conn
		cDone     <-chan struct{}
		sDone     <-chan struct{}
		handlerCh chan interface{}
		err       error
		onMessage onMessageFunc
		codec     Codec
		handler   Handler
	)

	sc := c.(*ServerConn)
	rawConn = sc.rawConn
	cDone = sc.ctx.Done()
	sDone = sc.belong.ctx.Done()
	handlerCh = sc.handlerCh
	codec = sc.belong.opts.codec
	handler = sc.belong.opts.hanlder
	onMessage = sc.belong.opts.onMessage

	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("%s panics: %v\n", sc.name, p)
		}
		wg.Done()
		xlogger.Debug("readLoop go-routine exited")
		c.Close()
	}()

	splitFunc := codec.GetScanSplitFun()
	scanner := bufio.NewScanner(rawConn)
	scanner.Split(splitFunc)

	for {
		select {
		case <-cDone: // connected closed
			xlogger.Debug(sc.name, ": read loop receiving cancel signal from conn")
			return
		case <-sDone: // server closed
			xlogger.Debug(sc.name, ": read loop receiving cancel signal from server")
			return
		default:

			if ok := scanner.Scan(); ok {
				//xlogger.Info("scanner msg ...")
				if msg, e := codec.Decode(scanner.Bytes()); e != nil {
					xlogger.Errorf("%s: error decoding message %v\n", sc.name, e)
					err = e
					return
				} else {
					//setHeartBeatFunc(time.Now().UnixNano())
					if handler == nil {
						if onMessage != nil {
							onMessage(msg, c.(WriteCloser))
						} else {
							xlogger.Warnf("readLoop no handler or onMessage() found for message\n")
						}
					}

					//xlogger.Info("put msg to channel ... ")

					handlerCh <- msg

					continue
				}

			} else {
				err = scanner.Err()
				if err != nil {
					xlogger.Errorf("%s: error scan bytes %v\n", sc.name, err)
					//if _, ok := err.(ErrUndefined); ok {
					//TODO update heart beats
					//setHeartBeatFunc(time.Now().UnixNano())
					//	continue
					//}
					return
				} else {
					// read data EOR
					xlogger.Infof("%s: read EOF.. \n",sc.name)
					c.Close()

				}

			}

		}
	}

}

func writeLoop(c WriteCloser, wg *sync.WaitGroup) {
	var (
		rawConn  net.Conn
		sendCh   chan []byte
		cDone    <-chan struct{}
		sDone    <-chan struct{}
		pkt      []byte
		err      error
		osLinger bool
	)

	sc := c.(*ServerConn)
	rawConn = sc.rawConn
	cDone = sc.ctx.Done()
	sDone = sc.belong.ctx.Done()
	// TODO 完善是否发送后续包判断 client , udp
	// 当为服务器链接时,serverConn在关闭后调用 OS_LINGER == 0
	osLinger = true
	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("%s: panics: %v\n", sc.name, p)
		}
		// 当未设置 OS_LINGER
		if !osLinger {
			// drain all pending messages before exit
			// 尝试将所有未发送完的数据发送完毕, 如果服务器端关闭将出错
		OuterFor:
			for {
				select {
				case pkt = <-sendCh:
					if pkt != nil {
						if _, err = rawConn.Write(pkt); err != nil {
							xlogger.Errorf("%s: error writing data %v\n", sc.name, err)
						}
					}
				default:
					break OuterFor
				}
			}
		}
		wg.Done()
		xlogger.Debugf("%s: writeLoop go-routine exited", sc.name)
		c.Close()

	}()

	// 循环发送数据包

	for {
		select {
		case <-cDone:
			xlogger.Debugf("%s: writeLoop receiving cancel signal from conn", sc.name)
			return
		case <-sDone:
			xlogger.Debugf("%s: writeLoop receiving cancel signal from server", sc.name)
			return
		case pkt = <-sendCh:
			if pkt != nil {
				if _, err = rawConn.Write(pkt); err != nil {
					xlogger.Errorf("%s: writeLoop error writing data %v\n", sc.name, err)
					return
				}
			}
		}
	}

}

func handleLoop(c WriteCloser, wg *sync.WaitGroup) {
	var (
		cDone <-chan struct{}
		sDone <-chan struct{}
		//timerCh      chan *OnTimeOut
		handlerCh chan interface{}
		//netID        int64
		//ctx          context.Context
		//askForWorker bool
		err     error
		hanlder Handler
	)
	sc := c.(*ServerConn)
	cDone = sc.ctx.Done()
	sDone = sc.belong.ctx.Done()
	//timerCh = c.timerCh
	handlerCh = sc.handlerCh
	//netID = sc.netId
	//ctx = sc.ctx
	hanlder = sc.belong.opts.hanlder
	//askForWorker = true

	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("%s: panics: %v\n", sc.name, p)
		}
		wg.Done()
		xlogger.Debugf("%s: handleLoop go-routine exited", sc.name)

		OuterFor:
		for {
			select {
			case msg := <-handlerCh:
				//xlogger.Debugf("%s: hanlde msg ... ")
				if hanlder != nil {
					err = hanlder.Handle(msg, c)
					if err != nil {
						xlogger.Errorf("%s: handleloop handle msg error:", sc.name, err)
					}
				}
			// TODO do some thing for timeout??
			default:
				break OuterFor
			}
		}
		c.Close()
	}()

	for {
		select {
		case <-cDone:
			xlogger.Debugf("%s: handleloop receiving cancel signal from conn", sc.name)
			return
		case <-sDone:
			xlogger.Debugf("%s: handleloop receiving cancel signal from server", sc.name)
			return
		case msg := <-handlerCh:
			xlogger.Debugf("%s: hanlde msg ... ")
			if hanlder != nil {
				err = hanlder.Handle(msg, c)
				if err != nil {
					xlogger.Errorf("%s: handleloop handle msg error:", sc.name, err)
				}
			}
			// TODO do some thing for timeout??
		}
	}

}
