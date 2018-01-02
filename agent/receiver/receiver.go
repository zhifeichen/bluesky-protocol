package receiver

import (
	"net"
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/agent/protocol/bluesky"
	"github.com/zhifeichen/bluesky-protocol/common/logger"
	config "github.com/zhifeichen/bluesky-protocol/agent/cfg"
)

// Start start receive
func Start() {
	ip := config.Config().Ip
	if ip == "" {
		logger.Error.Panicln("socket listen配置不正确", config.Config().Ip)
		return
	}
	port := config.Config().Port
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		logger.Fatal.Fatal("生成监听地址:", tcpAddr, " 失败:", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Info.Println("监听:", tcpAddr, " 失败:", err)
	}
	defer listener.Close()

	for{
		conn, err := listener.Accept()
		if err != nil{
			logger.Info.Println("获取客户端连接失败:",err)
			continue
		}
		fmt.Printf("client[%s] connected\n", conn.RemoteAddr().String())
		msgChan := make(chan []byte, 1)
		go bluesky.Get(conn, msgChan)
		i := 0
		for {
			msg := <- msgChan
			if len(msg) == 0 {
				//fmt.Println("msg is eof")
				break
			}
			fmt.Printf("receiced[%d]: %v; message is %v\n", i, msg, bluesky.CheckCRC(msg))
			i++
		}
	}
}