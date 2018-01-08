package server

import (
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"net"
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/common/tcpServer"
)

type BlueSkyProtocolServer struct {
	*tcpServer.Server
	ip   string
	port int
}


func NewBlueSkyProtocolServer(ip string, port int, handler tcpServer.Handler) *BlueSkyProtocolServer {
	onConnect := tcpServer.OnConnectOption(func(conn tcpServer.WriteCloser) bool {
		sc := conn.(*tcpServer.ServerConn)
		xlogger.Info("new blueSkyProtocol conn ", sc.RemoteAddr().String(), "...")
		return true
	})

	onDisConnect := tcpServer.OnCloseOption(func(conn tcpServer.WriteCloser) {
		sc := conn.(*tcpServer.ServerConn)
		xlogger.Info("blueSkyProtocol conn ", sc.RemoteAddr().String(), " disconnect.")
	})

	onCodec := tcpServer.OnCustomCodecOption(&BlueSkyCodec{})

	onHandler := tcpServer.CustomHandlerOption(handler)

	server, _ := tcpServer.NewServer(onConnect, onDisConnect, onCodec,onHandler)
	return &BlueSkyProtocolServer{
		server,
		ip,
		port,
	}
}

func StartTcpServer(ip string, port int, handler tcpServer.Handler) (error, *BlueSkyProtocolServer) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		xlogger.Error("start listen error", err)
		return err, nil
	}
	svr := NewBlueSkyProtocolServer(ip, port, handler)
	svr.Server.Start(l)
	return nil, svr
}
