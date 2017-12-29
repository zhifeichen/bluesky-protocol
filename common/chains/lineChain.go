package chains

import (
	"strings"
	"bluesky-protocol/common/logger"
	"bluesky-protocol/common/utils"
)

type LineChain struct {
	name string
	items []IItem
	msgs chan *ChainMsg
}


func NewLineChains(name string) *LineChain{
	return &LineChain{
		name,
		make([]IItem,0),
		make(chan *ChainMsg),
	}
}



func (c *LineChain) String() {
	names := make([]string, 0)
	for _, v := range c.items {
		names = append(names, v.GetName())
	}
	logger.Info.Println(" 线性处理链: ",c.name," task:", strings.Join(names, " -> "))
}

func (c *LineChain)AddItem(items ...IItem) error{
	msg := NewAddItemMsg(items,false)
	err,_ := c.HandleMsg(msg)
	return err
}

func (c *LineChain)HandleMsg(msg *ChainMsg) (error,interface{}){
	c.msgs <- msg
	if msg.sync && msg.syncChan != nil{
		if msg,ok := <- msg.syncChan;!ok{
			return common.NewError(common.CHAIN_HANDLE_MSG_ERROR),nil
		} else {
			return nil,msg.d
		}

	} else {
		return nil,nil
	}
}