package receiver

import (
	"fmt"
	config "github.com/zhifeichen/bluesky-protocol/agent/cfg"
	"github.com/zhifeichen/bluesky-protocol/common/logger"
	"github.com/zhifeichen/bluesky-protocol/common/protocol/bluesky"
	"net"
)

// Start start receive
func Start() {
	ip := config.Config().Ip
	if ip == "" {
		logger.Error("socket listen配置不正确", config.Config().Ip)
		return
	}
	port := config.Config().Port
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		logger.Fatal("生成监听地址:", tcpAddr, " 失败:", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Info("监听:", tcpAddr, " 失败:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Info("获取客户端连接失败:", err)
			continue
		}
		fmt.Printf("client[%s] connected\n", conn.RemoteAddr().String())
		go func(conn net.Conn) {
			msgChan := make(chan []byte, 1)
			go bluesky.Get(conn, msgChan)
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
			conn.Close()
		}(conn)
	}
}
