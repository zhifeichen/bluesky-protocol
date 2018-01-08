package servers

import (
	"github.com/zhifeichen/bluesky-protocol/common/tcpServer"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
)

func BlueSkyHandler(msg interface{},c tcpServer.WriteCloser) error{
	// TODO handler protocol pkt
	xlogger.Debugf("not implement msg...")
	return nil
}