package server

import (
	"time"
	"fmt"
	"net"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"sync"
	"bufio"
	"io"
)

type writePackage struct {
	msg []byte
	writer io.Writer
}

type readPackage struct {
	msg interface{}
	id uint64
	remote net.Addr
}

type processer struct {
	msgChn chan readPackage
	sendChn chan writePackage
}

type collectorInfo struct {
	prot ProtocolHandler
	remote net.Addr
}

type writer struct {
	w io.Writer
	timerout *time.Timer
}

var (
	proc = processer{msgChn: make(chan readPackage, 100), sendChn: make(chan writePackage, 100)}
	lock = sync.RWMutex{}
	writerMap = sync.Map{}
	collecterInfoMap = sync.Map{}
)

func init() {
	// timeoutDuration = time.Duration(cfg.Config().Timeout) * time.Second
	go writeLoop()
	go handleLoop()
}

func getWriterKey(addr net.Addr) string {
	n := addr.Network()
	s := addr.String()
	ret := n + "_" + s
	// xlogger.Debugf("get writer key: %s\n", ret)
	return ret
}

func readLoop(r io.ReadCloser, prot ProtocolHandler, remote net.Addr) {
	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("recover error: %v", p)
		}
		r.Close()
	}()
	scanner := bufio.NewScanner(r)
	scanner.Split(prot.Split)
	for scanner.Scan() {
		msgBuf := scanner.Bytes()
		if msg, id, err := prot.Decode(msgBuf); err != nil {
			xlogger.Errorf("decode msg[%v] error: %v", msgBuf, err)
		} else {
			if info, ok := collecterInfoMap.Load(id); !ok {
				collecterInfoMap.Store(id, collectorInfo{prot, remote})
			} else {
				i := info.(collectorInfo)
				oldKey := getWriterKey(i.remote)
				key := getWriterKey(remote)
				if oldKey == key {
					if w, ok := writerMap.Load(key); ok {
						t := w.(writer).timerout
						if !t.Stop() {
							<-t.C
						}
						t.Reset(timeoutDuration)
					}
				} else {
					collecterInfoMap.Store(id, collectorInfo{prot, remote})
				}
			}
			proc.msgChn <- readPackage{msg, id, remote}
		}
	}
	// key := getWriterKey(remote)
	// if writer, ok := writerMap[key]; ok {
	// 	c := atomic.AddInt32(&writer.count, -1)
	// 	if c == 0 {
	// 		// delete(writerMap, key)
	// 	}
	// }
	// xlogger.Debugf("colletror info map len: %d, writer map len: %d\n", len(collecterInfoMap), len(writerMap))
}

func writeLoop() {
	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("recover error: %v", p)
		}
	}()

	for {
		msg, ok := <-proc.sendChn
		if !ok {
			break
		}
		_, err := msg.writer.Write(msg.msg)
		if err != nil {
			xlogger.Errorf("send msg[%v] error: %v\n", msg, err)
		}
	}
}

func handleLoop() {
	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("recover error: %v", p)
		}
	}()

	for {
		msg, ok := <-proc.msgChn
		if !ok {
			break
		}
		if c, ok := collecterInfoMap.Load(msg.id); ok {
			c.(collectorInfo).prot.Handle(msg.msg, msg.remote)
		}
	}
}

func sendMsg(msg interface{}, info collectorInfo) error {
	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("recover error: %v", p)
		}
	}()

	msgBuf, err := info.prot.Encode(msg)
	if err != nil {
		xlogger.Errorf("encode msg error: %v\n", err)
		return err
	}
	if w, ok := writerMap.Load(getWriterKey(info.remote)); ok {
		proc.sendChn <- writePackage{msgBuf, w.(writer).w}
		return nil
	}
	return fmt.Errorf("remote[%s %s] is offline", info.remote.Network(), info.remote.String())
}

func Process(r io.ReadCloser, w io.Writer, handler ProtocolHandler, remote net.Addr) {
	key := getWriterKey(remote)
	if ow, ok := writerMap.Load(key); !ok {
		wr := writer{w, time.AfterFunc(timeoutDuration, func(){
			writerMap.Delete(key)
			xlogger.Debugf("colletror info map: %+v, writer map: %d\n", collecterInfoMap, writerMap)
		})}
		writerMap.Store(key, wr)
	} else {
		oww := ow.(writer)
		if !oww.timerout.Stop() {
			<- oww.timerout.C
		}
		oww.timerout.Reset(timeoutDuration)
	}
	go readLoop(r, handler, remote)
	// xlogger.Debugf("writeMap: %+v\n", writerMap)
}

func Stop() {
	close(proc.msgChn)
	close(proc.sendChn)
}

func SendMsg(id uint64, msg interface{}) {
	if info, ok := collecterInfoMap.Load(id); ok {
		sendMsg(msg, info.(collectorInfo))
	}
}