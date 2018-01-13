package server

import (
	"github.com/zhifeichen/bluesky-protocol/agent/cfg"
	"time"
	"io/ioutil"
	"bytes"
	"errors"
	"fmt"
	"net"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
)

type protocol interface {
	Split(data []byte, atEOF bool) (advance int, token []byte, err error)
	Decode(msg []byte) (interface{}, uint64, error)
	Encode(msg interface{}) ([]byte, error)
}

type handler interface {
	Handle(msg interface{}, remote net.Addr)
}

type ProtocolHandler interface {
	protocol
	handler
}

var (
	timeoutDuration time.Duration
)

func Init() {
	timeoutDuration = time.Duration(cfg.Config().Timeout) * time.Second
}

// TCPServer type
type TCPServer struct {
	server *net.TCPListener
	ProtocolHandler
}

func NewTCPServer(ip string, port int, protocolHandler ProtocolHandler) (*TCPServer, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		xlogger.Fatal("生成TCP监听地址:", tcpAddr, " 失败:", err)
		return nil, err
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		xlogger.Error("start tcp listen error", err)
		return nil, err
	}
	server := TCPServer{l, protocolHandler}
	return &server, nil
}

func (s *TCPServer)Start() {
	for {
		conn, err := s.server.AcceptTCP()
		if err != nil {
			xlogger.Error("获取客户端连接失败: ", err)
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			} else {
				break
			}
		}
		xlogger.Infof("client[%s] connnected.\n", conn.RemoteAddr().String())
		go Process(conn, conn, s.ProtocolHandler, conn.RemoteAddr())
	}
	s.server.Close()
}

func (s *TCPServer)Stop() {
	s.server.Close()
}

type udpWrite struct {
	w *net.UDPConn
	addr net.Addr
}

func (w udpWrite) Write(p []byte) (n int, err error) {
	return w.w.WriteTo(p, w.addr)
}

// UDPServer type
type UDPServer struct {
	server *net.UDPConn
	ProtocolHandler
}

func NewUDPServer(ip string, port int, protocolHandler ProtocolHandler) (*UDPServer, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		xlogger.Fatal("生成udp监听地址:", udpAddr, " 失败:", err)
		return nil, errors.New("udp地址错误")
	}
	l, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		xlogger.Error("start udp listen error", err)
		return nil, err
	}
	server := UDPServer{l, protocolHandler}
	return &server, nil
}

func (s *UDPServer)Start() {
	const MAXUDPPACKLEN int = 1500
	for {
		buffer := make([]byte, MAXUDPPACKLEN)
		n, remote, err := s.server.ReadFrom(buffer)
		if err != nil {
			xlogger.Error("udp 读取数据失败: ", err)
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			} else {
				break
			}
		}
		go Process(ioutil.NopCloser(bytes.NewReader(buffer[:n])), udpWrite{s.server, remote}, s.ProtocolHandler, remote)
	}
	s.server.Close()
}

func (s *UDPServer)Stop() {
	s.server.Close()
}
