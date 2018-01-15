package pingPong

import (
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"github.com/zhifeichen/bluesky-protocol/common/server"
	"net"
)

type PingPongHandler struct {
}

func (c *PingPongHandler) Handle(msg interface{}, wc server.WriteCloser) error {
	xlogger.Debug("handler tcp msg:", msg)
	wc.WriteTCP(msg)
	return nil
}
func (c *PingPongHandler) HandleUdp(msg interface{}, wc server.WriteCloser,
	remoteAddr *net.UDPAddr) error {
	xlogger.Debug("handler udp msg:", msg)
	wc.WriteUDP(msg, remoteAddr)
	return nil
}
