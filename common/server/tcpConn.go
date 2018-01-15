package server

import (
	"bufio"
	"context"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"net"
	"sync"
)

// TCP server conn
type TcpServerConn struct {
	ServerConn
}

func NewTcpServerConn(id int64, s *TCPServer, c net.Conn) *TcpServerConn {
	sc := &TcpServerConn{
		ServerConn: ServerConn{
			netId:   id,
			belong:  &s.server,
			rawConn: c,
			once:    &sync.Once{},
			wg:      &sync.WaitGroup{},
			sendCh:    make(chan connSendMsg, s.opts.bufferSize),
			handlerCh: make(chan connHandleMsg, s.opts.bufferSize),
		},

	}

	sc.ctx, sc.cancel = context.WithCancel(context.WithValue(s.ctx, serverCtx, s))

	sc.name = "tcp_" + c.RemoteAddr().String()

	return sc
}

// Start starts the server connection, creating go-routines for reading,
// writing and handlng.
func (sc *TcpServerConn) Start() {
	xlogger.Infof("conn start, <%v -> %v>\n", sc.rawConn.LocalAddr(), sc.rawConn.RemoteAddr())
	onConnect := sc.belong.opts.onConnect
	if onConnect != nil {
		onConnect(sc)
	}

	loopers := []func(*TcpServerConn, *sync.WaitGroup){readLoop, writeLoop, handleLoop}

	for _, l := range loopers {
		sc.wg.Add(1)
		go l(sc, sc.wg)
	}

}

func (sc *TcpServerConn) Write(msg interface{}) error {
	xlogger.Warnf("%s: writeUDP ? maybe use WriteTCP(msg interface{}) instand!")
	return asyncWrite(sc, msg)

}

func asyncWrite(sc *TcpServerConn, msg interface{}) (err error) {
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
		xlogger.Errorf("asyncWrite error %v\n", err)
		return
	}

	//xlogger.Debug("asyncWrite:", msg)

	select {
	case sendCh <- connSendMsg{pkt, sc.Addr()}:
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
func readLoop(sc *TcpServerConn, wg *sync.WaitGroup) {

	var (
		rawConn   net.Conn
		cDone     <-chan struct{}
		sDone     <-chan struct{}
		handlerCh chan connHandleMsg
		err       error
		onMessage onMessageFunc
		codec     Codec
		handler   Handler
	)

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
		sc.Close()
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
							onMessage(msg, sc)
						} else {
							xlogger.Warnf("readLoop no handler or onMessage() found for message\n")
						}
					}

					//xlogger.Info("put msg to channel ... ")

					handlerCh <- connHandleMsg{d: msg, removeAddr: sc.Addr()}

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
					xlogger.Infof("%s: read EOF.. \n", sc.name)
					//c.Close()
					return

				}

			}

		}
	}

}

func writeLoop(sc *TcpServerConn, wg *sync.WaitGroup) {
	var (
		rawConn  net.Conn
		sendCh   chan connSendMsg
		cDone    <-chan struct{}
		sDone    <-chan struct{}
		pkt      connSendMsg
		err      error
		osLinger bool
	)

	rawConn = sc.rawConn
	cDone = sc.ctx.Done()
	sDone = sc.belong.ctx.Done()
	sendCh = sc.sendCh
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
					xlogger.Debugf("%s: 发送数据:", sc.name, pkt)
					if pkt.d != nil {
						if _, err = rawConn.Write(pkt.d); err != nil {
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
		sc.Close()

	}()
	xlogger.Debugf("%s: writeLoop go-routine start ...", sc.name)
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
			//xlogger.Debugf("%s: send msg:%v ...", sc.name, pkt)
			if pkt.d != nil {
				if _, err = rawConn.Write(pkt.d); err != nil {
					xlogger.Errorf("%s: writeLoop error writing data %v\n", sc.name, err)
					return
				}
			}
		}
	}

}

func handleLoop(sc *TcpServerConn, wg *sync.WaitGroup) {
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

		xlogger.Debugf("%s: handleLoop go-routine exited", sc.name)

	OuterFor:
		for {
			select {
			case msg := <-handlerCh:
				//xlogger.Debugf("%s: hanlde msg ... ")
				if hanlder != nil {
					err = hanlder.Handle(msg.d, sc)
					if err != nil {
						xlogger.Errorf("%s: handleloop handle msg error:", sc.name, err)
					}
				}
				// TODO do some thing for timeout??
			default:
				break OuterFor
			}
		}
		wg.Done()
		sc.Close()
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
			if hanlder != nil {
				err = hanlder.Handle(msg.d, sc)
				if err != nil {
					xlogger.Errorf("%s: handleloop handle msg error:", sc.name, err)
				}
			}
			// TODO do some thing for timeout??
		}
	}

}
