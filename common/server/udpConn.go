package server

import (
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"net"
	"sync"
	"context"
	"bytes"
	"bufio"
	"errors"
)

// udp server
type udpSendMsg struct {
	d          []byte
	removeAddr *net.UDPAddr
}

type udpHandleMsg struct {
	d          interface{}
	removeAddr *net.UDPAddr
}

type udpMsgConn struct {
	udpServerConn *UdpServerConn
	removeAddr    *net.UDPAddr
}

func (c *udpMsgConn) Write(msg interface{}) error {
	return asyncUdpWrite(c.udpServerConn, msg, c.removeAddr)
}
func (c *udpMsgConn) Close() {
	xlogger.Warnf("%s: you should never close a udp conn from server!")
	return
}

type UdpServerConn struct {
	mtu     int
	netId   int64
	belong  *UDPServer
	rawConn *net.UDPConn
	name    string

	once      *sync.Once
	wg        *sync.WaitGroup
	sendCh    chan udpSendMsg
	handlerCh chan udpHandleMsg
	//timerCh   chan *OnTimeOut

	mu     sync.Mutex // guards following
	ctx    context.Context
	cancel context.CancelFunc
}

func NewUdpServerConn(id int64, s *UDPServer, c *net.UDPConn) *UdpServerConn {
	sc := &UdpServerConn{
		mtu:       1500,
		netId:     id,
		belong:    s,
		rawConn:   c,
		once:      &sync.Once{},
		wg:        &sync.WaitGroup{},
		sendCh:    make(chan udpSendMsg, s.opts.bufferSize),
		handlerCh: make(chan udpHandleMsg, s.opts.bufferSize),
	}

	sc.ctx, sc.cancel = context.WithCancel(context.WithValue(s.ctx, serverCtx, s))

	sc.name = "udp_" + c.LocalAddr().String()

	return sc
}

func (sc *UdpServerConn) SetName(name string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.name = name
}

func (sc *UdpServerConn) GetName() string {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.name
}

// Start starts the server connection, creating go-routines for reading,
// writing and handlng.
func (sc *UdpServerConn) Start() {
	xlogger.Infof("udp conn start, %v", sc.rawConn.LocalAddr())
	onConnect := sc.belong.opts.onConnect
	if onConnect != nil {
		onConnect(sc)
	}

	loopers := []func(WriteCloser, *sync.WaitGroup){udpReadLoop, udpWriteLoop, handleUdpLoop}

	for _, l := range loopers {
		sc.wg.Add(1)
		go l(sc, sc.wg)
	}

}

// Close gracefully closes the server connection. It blocked until all sub
// go-routines are completed and returned.
func (sc *UdpServerConn) Close() {
	// TODO??
	sc.once.Do(func() {
		xlogger.Infof("udp conn close gracefully, <%v -> %v>\n", sc.rawConn.LocalAddr())

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

func (sc *UdpServerConn) Write(msg interface{}) error {
	return errors.New("you should never write a udp server as tcp")
}

func asyncUdpWrite(sc *UdpServerConn, msg interface{}, remoteAddr *net.UDPAddr) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = ErrServerClosed
		}
	}()

	var (
		pkt    []byte
		sendCh chan udpSendMsg
	)
	pkt, err = sc.belong.opts.codec.Encode(msg)
	sendCh = sc.sendCh

	if err != nil {
		xlogger.Errorf("udp asyncWrite error %v\n", err)
		return
	}

	select {
	case sendCh <- udpSendMsg{d: pkt, removeAddr: remoteAddr}:
		err = nil
	default:
		err = ErrWouldBlock
	}
	return
}

/**
*/
func udpReadLoop(c WriteCloser, wg *sync.WaitGroup) {

	var (
		cDone <-chan struct{}
		sDone <-chan struct{}
	)

	sc := c.(*UdpServerConn)

	cDone = sc.ctx.Done()
	sDone = sc.belong.ctx.Done()
	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("%s panics: %v\n", sc.name, p)
		}
		wg.Done()
		xlogger.Debug("readLoop go-routine exited")
		sc.Close()
	}()

	for {

		select {
		case <-cDone: // connected closed
			xlogger.Debug(sc.name, ": read loop receiving cancel signal from conn")
			return
		case <-sDone: // server closed
			xlogger.Debug(sc.name, ": read loop receiving cancel signal from server")
			return
		default:
			doReadLoop(sc)
		}
	}
}

func doReadLoop(sc *UdpServerConn) {
	var (
		rawConn   *net.UDPConn
		handlerCh chan udpHandleMsg
		onMessage onMessageFunc
		codec     Codec
		handler   Handler
	)
	rawConn = sc.rawConn
	codec = sc.belong.opts.codec
	handler = sc.belong.opts.hanlder
	onMessage = sc.belong.opts.onMessage
	handlerCh = sc.handlerCh

	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("%s read panics: %v\n", sc.name, p)
		}
	}()

	buffer := make([]byte, sc.mtu)
	if n, remote, err := rawConn.ReadFromUDP(buffer); err != nil {
		xlogger.Errorf("%s: read data from %v failed: %v\n", sc.name, remote, err)
	} else {
		xlogger.Debug("read data from udp ... n:", n, " err:", err)
		scanner := bufio.NewScanner(bytes.NewReader(buffer[:n]))
		scanner.Split(codec.GetScanSplitFun())
		if ok := scanner.Scan(); ok {
			if msg, e := codec.Decode(scanner.Bytes()); e != nil {
				xlogger.Errorf("%s: error decoding message %v\n", sc.name, e)
			} else {
				//setHeartBeatFunc(time.Now().UnixNano())
				if handler == nil {
					if onMessage != nil {
						onMessage(msg, sc)
					} else {
						xlogger.Warnf("%s readLoop no handler or onMessage() found for message\n", sc.name)
					}
				}
				xlogger.Info("put msg to channel ... ", msg)
				handlerCh <- udpHandleMsg{d: msg, removeAddr: remote}
			}
		} else {
			err = scanner.Err()
			if err != nil {
				xlogger.Errorf("%s: error scan bytes %v\n", sc.name, err)
				return
			} else {
				xlogger.Infof("%s: read EOF.. \n", sc.name)
				return

			}
		}
	}
}

func udpWriteLoop(c WriteCloser, wg *sync.WaitGroup) {
	var (
		rawConn *net.UDPConn
		sendCh  chan udpSendMsg
		cDone   <-chan struct{}
		sDone   <-chan struct{}
		msg     udpSendMsg
		err     error
	)

	sc := c.(*UdpServerConn)
	rawConn = sc.rawConn
	cDone = sc.ctx.Done()
	sDone = sc.belong.ctx.Done()
	sendCh = sc.sendCh
	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("%s: write panics: %v\n", sc.name, p)
		}
		wg.Done()
		xlogger.Debugf("%s: writeLoop go-routine exited", sc.name)
		sc.Close()
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
		case msg = <-sendCh:
			//xlogger.Debugf("%s: write msg:%v",sc.name,msg.d,msg.removeAddr)
			if _, err = rawConn.WriteToUDP(msg.d, msg.removeAddr); err != nil {
				xlogger.Errorf("%s: writeLoop error writing data %v\n", sc.name, err)
			}
		}
	}

}

func handleUdpLoop(c WriteCloser, wg *sync.WaitGroup) {
	var (
		cDone <-chan struct{}
		sDone <-chan struct{}
		//timerCh      chan *OnTimeOut
		handlerCh chan udpHandleMsg
		//netID        int64
		//ctx          context.Context
		//askForWorker bool
		err     error
		hanlder Handler
	)
	sc := c.(*UdpServerConn)
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
		sc.Close()
	}()

	for {
		//xlogger.Debugf("%s: hanlde msg loop ... ")
		select {
		case <-cDone:
			xlogger.Debugf("%s: handleloop receiving cancel signal from conn", sc.name)
			return
		case <-sDone:
			xlogger.Debugf("%s: handleloop receiving cancel signal from server", sc.name)
			return
		case msg := <-handlerCh:
			//xlogger.Debugf("%s: hanlde msg ... ",sc.name)
			if hanlder != nil {
				err = doHandleLoop(hanlder, msg, sc)
				if err != nil {
					xlogger.Errorf("%s: handleloop handle msg error:", sc.name, err)
				}
			}
			// TODO do some thing for timeout??
		}
	}
}

func doHandleLoop(hanlder Handler, msg udpHandleMsg, sc *UdpServerConn) (err error) {
	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("%s handle panics: %v\n", sc.name, p)
		}
	}()
	err = hanlder.Handle(msg.d, &udpMsgConn{sc, msg.removeAddr})
	return
}
