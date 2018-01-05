package test

import (
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("测试 日志", func(t *testing.T) {
		xlogger.Info("test")
		xlogger.New("./logs/access.log", xlogger.DEBUG, true)
		defer xlogger.Close()
		xlogger.Info("test")
		xlogger.Warn("test")
		xlogger.Error("test")

	})

}
