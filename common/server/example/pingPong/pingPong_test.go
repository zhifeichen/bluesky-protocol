package pingPong

import (
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"net"
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/common/server"
	"errors"
	"testing"
	"time"
)

const (
	max_msg int = 20
	max_gos int = 10
)

type PingPongServer struct {
	*server.TCPServer
	ip   string
	port int
}

func NewPingPongServer(ip string, port int, handler server.Handler) *PingPongServer {
	onConnect := server.OnConnectOption(func(conn server.WriteCloser) bool {
		sc := conn.(*server.TcpServerConn)
		xlogger.Info("new pingpongprotocol conn ", sc.Addr().String(), "...")
		return true
	})

	onDisConnect := server.OnCloseOption(func(conn server.Closer) {
		sc := conn.(*server.ServerConn)
		xlogger.Info("pingpongprotocol conn ", sc.Addr().String(), " disconnect.")
	})

	onCodec := server.OnCustomCodecOption(&PinPongCodec{})

	onHandler := server.CustomHandlerOption(handler)

	server, _ := server.NewTCPServer(onConnect, onDisConnect, onCodec, onHandler)
	return &PingPongServer{
		server,
		ip,
		port,
	}
}

func StartTcpServer(ip string, port int, handler server.Handler) (error, *PingPongServer) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		xlogger.Error("start listen error", err)
		return err, nil
	}
	svr := NewPingPongServer(ip, port, handler)
	go svr.TCPServer.Start(l)
	return nil, svr
}

type PingPongUdpServer struct {
	*server.UDPServer
	ip   string
	port int
}

func NewPingPongUdpServer(ip string, port int, handler server.Handler) *PingPongUdpServer {

	onCodec := server.OnCustomCodecOption(&PinPongCodec{})

	onHandler := server.CustomHandlerOption(handler)

	server, _ := server.NewUDPServer(onCodec, onHandler)
	return &PingPongUdpServer{
		server,
		ip,
		port,
	}
}

func StartUdpServer(ip string, port int, handler server.Handler) (error, *PingPongUdpServer) {
	//xlogger.Debugf("建立udp侦听: %s:%d",ip,port)
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		xlogger.Fatal("生成udp监听地址:", udpAddr, " 失败:", err)
		return errors.New("udp地址错误"), nil
	}
	l, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		xlogger.Error("start udp listen error", err)
		return err, nil
	}
	svr := NewPingPongUdpServer(ip, port, handler)
	go svr.UDPServer.Start(l)
	return nil, svr
}

func TestPingPong(t *testing.T) {
	xlogger.New("./logs/x.log", xlogger.DEBUG, true)
	defer xlogger.Close()
	t.Run("测试 chain", func(t *testing.T) {
		doneChan := make(chan int, max_msg)
		handler := &PingPongHandler{}
		ip, port := "127.0.0.1", 8090
		_, svr := StartTcpServer(ip, port, handler)
		_, svrUdp := StartUdpServer(ip, port, handler)
		//go StartUdpServer(ip, port, handler)
		time.Sleep(time.Second)
		for i := 0; i < max_gos; i++ {
			go StartTcpClientSendMsg(ip, port, max_msg, doneChan)
			go StartUdpClientSendMsg(ip, port, max_msg, doneChan)
		}

		//fmt.Println("wait msg is done...")
		var i int
		timeout := time.NewTimer(time.Second * 20)
	OuterFor:
		for {
			select {
			case <-doneChan:
				i++
				if i >= max_msg*max_gos*2 {
					break OuterFor
				}
			case <-timeout.C:
				fmt.Println("测试超时")
				t.Fatal()
			}
		}

		fmt.Println("测试完成 ... [ok]")
		svr.TCPServer.Stop()
		svrUdp.UDPServer.Stop()
	})

}
