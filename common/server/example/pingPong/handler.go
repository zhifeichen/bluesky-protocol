package pingPong

import (
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"github.com/zhifeichen/bluesky-protocol/common/server"
)

type PingPongHandler struct {
}

func (c *PingPongHandler) Handle(msg interface{}, wc server.WriteCloser) error {
	xlogger.Debug("handler tcp msg:", msg)
	wc.Write(msg)
	return nil
}
