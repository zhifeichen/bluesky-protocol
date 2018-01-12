package servers

import (
	"github.com/zhifeichen/bluesky-protocol/common/tcpServer"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"net"
	"errors"
)

type BlueSkyHandler struct {

}

func (h* BlueSkyHandler)Handle(msg interface{},c tcpServer.WriteCloser) error{
	// TODO handler protocol pkt
	xlogger.Debugf("not implement tcp  msg...")

	// TODO write ack ..
	//c.WriteTCP(xx)
	return nil
}

func (h* BlueSkyHandler)HandleUdp(msg interface{},c tcpServer.WriteCloser,remoteAddr *net.UDPAddr) error{
	// TODO handler protocol pkt
	xlogger.Debugf("not implement udp  msg...")

	// TODO write ack ..
	//c.WriteUDP(xx,remoteAddr)
	panic(errors.New("测试错误"))
	return errors.New("测试错误")
}