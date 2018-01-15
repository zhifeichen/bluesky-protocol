package server

import (
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"net"
	"sync"
	"context"
	"bytes"
	"bufio"
)

// udp server
type udpServerConnWrap struct {
	udpServerConn *UdpServerConn
	removeAddr    net.Addr
}

func (c *udpServerConnWrap) Write(msg interface{}) error {
	return asyncUdpWrite(c.udpServerConn, msg, c.removeAddr)
}
func (c *udpServerConnWrap) Close() {
	xlogger.Warnf("%s: you should never close a udp conn from server!")
	return
}

type UdpServerConn struct {
	ServerConn
}

func NewUdpServerConn(id int64, s *UDPServer, c *net.UDPConn) *UdpServerConn {
	sc := &UdpServerConn{
		ServerConn: ServerConn{netId: id,
			belong: &s.server,
			rawConn: c,
			once: &sync.Once{},
			wg: &sync.WaitGroup{},
			sendCh: make(chan connSendMsg, s.opts.bufferSize),
			handlerCh: make(chan connHandleMsg, s.opts.bufferSize),
		},
	}

	sc.ctx, sc.cancel = context.WithCancel(context.WithValue(s.ctx, serverCtx, s))

	sc.name = "udp_" + c.LocalAddr().String()

	return sc
}

// Start starts the server connection, creating go-routines for reading,
// writing and handlng.
func (sc *UdpServerConn) Start() {
	xlogger.Infof("udp conn start, %v", sc.rawConn.LocalAddr())
	loopers := []func(*UdpServerConn, *sync.WaitGroup){udpReadLoop, udpWriteLoop, handleUdpLoop}

	for _, l := range loopers {
		sc.wg.Add(1)
		go l(sc, sc.wg)
	}

}

// Addr returns the remote address of server connection.
func (sc *UdpServerConn) Addr() net.Addr {
	return sc.rawConn.LocalAddr()
}

func asyncUdpWrite(sc *UdpServerConn, msg interface{}, remoteAddr net.Addr) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = ErrServerClosed
		}
	}()

	var (
		pkt    []byte
		sendCh chan connSendMsg
	)
	pkt, err = sc.belong.opts.codec.Encode(msg)
	sendCh = sc.sendCh

	if err != nil {
		xlogger.Errorf("udp asyncWrite error %v\n", err)
		return
	}

	select {
	case sendCh <- connSendMsg{d: pkt, removeAddr: remoteAddr}:
		err = nil
	default:
		err = ErrWouldBlock
	}
	return
}

/**
*/
func udpReadLoop(sc *UdpServerConn, wg *sync.WaitGroup) {

	var (
		cDone <-chan struct{}
		sDone <-chan struct{}
	)

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
		handlerCh chan connHandleMsg
		onMessage onMessageFunc
		codec     Codec
		handler   Handler
	)
	rawConn = sc.rawConn.(*net.UDPConn)
	codec = sc.belong.opts.codec
	handler = sc.belong.opts.hanlder
	onMessage = sc.belong.opts.onMessage
	handlerCh = sc.handlerCh

	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("%s read panics: %v\n", sc.name, p)
		}
	}()

	buffer := make([]byte, 1500)
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
						onMessage(msg, &udpServerConnWrap{sc, remote})
					} else {
						xlogger.Warnf("%s readLoop no handler or onMessage() found for message\n", sc.name)
					}
				}
				xlogger.Info("put msg to channel ... ", msg)
				handlerCh <- connHandleMsg{d: msg, removeAddr: remote}
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

func udpWriteLoop(sc *UdpServerConn, wg *sync.WaitGroup) {
	var (
		rawConn *net.UDPConn
		sendCh  chan connSendMsg
		cDone   <-chan struct{}
		sDone   <-chan struct{}
		msg     connSendMsg
		err     error
	)

	rawConn = sc.rawConn.(*net.UDPConn)
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
			if _, err = rawConn.WriteToUDP(msg.d, msg.removeAddr.(*net.UDPAddr)); err != nil {
				xlogger.Errorf("%s: writeLoop error writing data %v\n", sc.name, err)
			}
		}
	}

}

func handleUdpLoop(sc *UdpServerConn, wg *sync.WaitGroup) {
	var (
		cDone <-chan struct{}
		sDone <-chan struct{}
		//timerCh      chan *OnTimeOut
		handlerCh chan connHandleMsg
		//netID        int64
		//ctx          context.Context
		//askForWorker bool
		err     error
		hanlder Handler
	)
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

func doHandleLoop(hanlder Handler, msg connHandleMsg, sc *UdpServerConn) (err error) {
	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("%s handle panics: %v\n", sc.name, p)
		}
	}()
	err = hanlder.Handle(msg.d, &udpServerConnWrap{sc, msg.removeAddr})
	return
}
