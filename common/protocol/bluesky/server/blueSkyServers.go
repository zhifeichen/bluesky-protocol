package server

import (
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"net"
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/common/server"
	"errors"
)

type BlueSkyProtocolServer struct {
	*server.TCPServer
	ip   string
	port int
}


func NewBlueSkyProtocolServer(ip string, port int, handler server.Handler) *BlueSkyProtocolServer {
	onConnect := server.OnConnectOption(func(conn server.WriteCloser) bool {
		sc := conn.(*server.ServerConn)
		xlogger.Info("new blueSkyProtocol conn ", sc.RemoteAddr().String(), "...")
		return true
	})

	onDisConnect := server.OnCloseOption(func(conn server.WriteCloser) {
		sc := conn.(*server.ServerConn)
		xlogger.Info("blueSkyProtocol conn ", sc.RemoteAddr().String(), " disconnect.")
	})

	onCodec := server.OnCustomCodecOption(&BlueSkyCodec{})

	onHandler := server.CustomHandlerOption(handler)

	server, _ := server.NewTCPServer(onConnect, onDisConnect, onCodec,onHandler)
	return &BlueSkyProtocolServer{
		server,
		ip,
		port,
	}
}

func StartTcpServer(ip string, port int, handler server.Handler) (error, *BlueSkyProtocolServer) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		xlogger.Error("start listen error", err)
		return err, nil
	}
	svr := NewBlueSkyProtocolServer(ip, port, handler)
	svr.TCPServer.Start(l)
	return nil, svr
}


type BlueSkyProtocolUdpServer struct {
	*server.UDPServer
	ip   string
	port int
}

func NewBlueSkyProtocolUdpServer(ip string, port int, handler server.Handler) *BlueSkyProtocolUdpServer {

	onCodec := server.OnCustomCodecOption(&BlueSkyCodec{})

	onHandler := server.CustomHandlerOption(handler)

	server, _ := server.NewUDPServer(onCodec,onHandler)
	return &BlueSkyProtocolUdpServer{
		server,
		ip,
		port,
	}
}

func StartUdpServer(ip string, port int, handler server.Handler) (error, *BlueSkyProtocolUdpServer) {
	//xlogger.Debugf("建立udp侦听: %s:%d",ip,port)
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		xlogger.Fatal("生成udp监听地址:", udpAddr, " 失败:", err)
		return errors.New("udp地址错误"),nil
	}
	l, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		xlogger.Error("start udp listen error", err)
		return err, nil
	}
	svr := NewBlueSkyProtocolUdpServer(ip, port, handler)
	svr.UDPServer.Start(l)
	return nil, svr
}
