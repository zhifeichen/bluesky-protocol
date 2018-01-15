package server

import (
	"net"
	"sync"
	"context"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
)

type Closer interface{
	Close()
}

// WriteCloser is the interface that groups Write and Close methods.
type WriteCloser interface {
	Write(interface{}) error
	Close()
}


type ServerConn struct {
	netId   int64
	belong  *server
	rawConn net.Conn
	name    string

	once      *sync.Once
	wg        *sync.WaitGroup
	//timerCh   chan *OnTimeOut

	mu     sync.Mutex // guards following
	ctx    context.Context
	cancel context.CancelFunc
	sendCh    chan connSendMsg
	handlerCh chan connHandleMsg
}

func (sc *ServerConn) Addr() net.Addr {
	return sc.rawConn.RemoteAddr()
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
		xlogger.Infof("%s: closed [ok]", sc.name)
		sc.belong.wg.Done()
	})
}


type connHandleMsg struct {
	d          interface{}
	removeAddr net.Addr
}


type connSendMsg struct {
	d          []byte
	removeAddr net.Addr
}
