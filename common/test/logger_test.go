package test

import (
	"testing"
	"github.com/zhifeichen/bluesky-protocol/common/logger"
)

func TestLogger(t *testing.T) {
	t.Run("测试 日志", func(t *testing.T) {
		logger.Info("test")
		logger.New("./logs/access.log",logger.DEBUG,true)
		defer logger.Close()
		logger.Info("test")
		logger.Warn("test")
		logger.Error("test")

	})

}