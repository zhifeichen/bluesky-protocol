package xlogger
import (
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("测试 日志", func(t *testing.T) {
		Info("test")
		New("./logs/access.log", DEBUG, true)
		defer Close()
		Info("test")
		Warn("test")
		Error("test")

	})

}
