package servers

import (
	"github.com/zhifeichen/bluesky-protocol/common/protocol/bluesky"
	"github.com/zhifeichen/bluesky-protocol/common/protocol/bluesky/server"
)

func Start(ip string, port int) {
	bluesky.RegisterAll()
	handler := &BlueSkyHandler{}
	go server.StartTcpServer(ip, port, handler)
	go server.StartUdpServer(ip, port, handler)

}
