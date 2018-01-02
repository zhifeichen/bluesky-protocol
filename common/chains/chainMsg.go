package chains

import (
	"fmt"
	"time"
)
/**
	追踪消息
 */
type ChainMsgTrace struct {
	Time     int64
	Duration int64
	Step     int8
	Ok       bool
}
/**
	chain消息
 */
type ChainMsg struct {
	Seqno    int64         			// 序号
	T        chainMsgType  			// 消息类型,main type
	Data     interface{}   			// data
	Sync     bool          			// 是否等待消息执行结果返回
	// TODO 如何更好的返回结果??
	syncChan chan *ChainMsgACK 			// 接收结果消息返回channel
	Track    bool					// 是否追踪消息
	Tracks   []ChainMsgTrace		// 追踪结果
}

func NewMsg(t chainMsgType, d interface{}, sync bool) *ChainMsg {
	msg := &ChainMsg{
		Seqno: time.Now().Unix(),
		T: t,
		Data: d,
		Sync:sync,
	}
	if sync {
		msg.syncChan = make(chan *ChainMsgACK)
	}
	return msg
}

func NewAddItemMsg(d interface{}, sync bool) *ChainMsg {
	return NewMsg(CHAIN_ADD_ITEM, d, sync)
}

func (c *ChainMsg)String() string {
	return fmt.Sprintf(
		"{seqno:%d, t:%v, sync:%v}",
		c.Seqno,
		c.T,
		c.Sync,
	)
}

type ChainMsgACK struct {
	Seqno int64        // 序号
	T     chainMsgType // 消息类型,main type
	Data  interface{}  // data
}

func NewMsgAck(seqno int64, t chainMsgType, d interface{}) *ChainMsgACK {
	msg := &ChainMsgACK{
		seqno,
		t,
		d,
	}
	return msg
}