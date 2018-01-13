package servers

import (
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"github.com/zhifeichen/bluesky-protocol/common/protocol/bluesky"
	"github.com/zhifeichen/bluesky-protocol/agent/servers/server"
)

func Start(ip string, port int) {
	bluesky.RegisterAll()
	handler := &BlueSkyHandler{}
	server.Init()
	tcpServer, err := server.NewTCPServer(ip, port, handler)
	if err != nil {
		xlogger.Errorf("start tcp server[%s:%d] error: %v\n", ip, port, err)
	}
	udpServer, err := server.NewUDPServer(ip, port, handler)
	if err != nil {
		xlogger.Errorf("start udp server[%s:%d] error: %v\n", ip, port, err)
	}
	go tcpServer.Start()
	go udpServer.Start()
}
