package servers

import (
	"github.com/zhifeichen/bluesky-protocol/common/protocol/bluesky"
	"github.com/zhifeichen/bluesky-protocol/common/protocol/bluesky/server"
	"github.com/zhifeichen/bluesky-protocol/common/tcpServer"
)

func Start(ip string, port int) {
	bluesky.RegisterAll()
	server.StartTcpServer(ip, port, tcpServer.HandlerFunc(BlueSkyHandler))
}
