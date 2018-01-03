package test

import (
	"testing"
	"github.com/zhifeichen/bluesky-protocol/common/chains"
	"fmt"
)

func TestChain(t *testing.T) {
	t.Run("测试 chain", func(t *testing.T) {
		lineChain := chains.NewLineChains("测试")
		fmt.Println(lineChain.String())
		lineChain.Run()
	})
}