package test

import (
	"testing"
	"github.com/zhifeichen/bluesky-protocol/common/chains"
	"fmt"
	"time"
)


type PrintItem struct {
	chains.BaseItem
	say string
}

func NewPrintItem(name,say string) *PrintItem{
	return &PrintItem{
		BaseItem:BaseItem{
			name:name,
		},
	}
}

func (item *PrintItem) Do(data interface{}) (error,interface{}){
	fmt.Println(item.GetName()," say:",item.say);
	return nil,data
}


func TestChain(t *testing.T) {
	t.Run("测试 chain", func(t *testing.T) {
		lineChain := chains.NewLineChains("测试")
		fmt.Println(lineChain.String())
		lineChain.Run()

		time.Sleep(time.Duration(1)*time.Second)
		//lineChain.S
	})
}