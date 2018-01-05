package receiver

import (
	"bytes"
	"fmt"
	config "github.com/zhifeichen/bluesky-protocol/agent/cfg"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"github.com/zhifeichen/bluesky-protocol/common/protocol/bluesky"
	"net"
)

func UdpStart() {
	addr := config.Config().UDPAddr
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		xlogger.Fatal("生成监听地址:", udpAddr, " 失败:", err)
	}
	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		xlogger.Info("监听:", udpAddr, " 失败:", err)
	}
	defer listener.Close()

	for {
		buffer := make([]byte, 1500)
		n, remote, err := listener.ReadFromUDP(buffer)
		if err != nil {
			xlogger.Error("读取 ", remote, " 数据失败")
			continue
		}
		if n < 30 {
			xlogger.Error("读取数据不完整, n: ", n)
			continue
		}
		msgChan := make(chan []byte, 1)
		go bluesky.Get(bytes.NewReader(buffer[:n]), msgChan)
		i := 0
		for {
			msg := <-msgChan
			if len(msg) == 0 {
				//fmt.Println("msg is eof")
				break
			}
			fmt.Printf("receiced[%d]: %v; message is %v\n", i, msg, bluesky.CheckCRC(msg))
			i++
			msgComm := bluesky.Common{}
			err := msgComm.Unmarshal(msg)
			if err != nil {
				continue
			}
			bluesky.HandleMessage(&msgComm)
		}
	}
}
