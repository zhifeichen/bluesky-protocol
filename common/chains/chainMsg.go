package chains

import (
	"fmt"
	"time"
	"context"
	"sync"
)
/**
	追踪消息
 */
type ChainTrace struct {
	Step     int
	Time     int64 // ms
	Duration int64
	Error    error
}

func (trace *ChainTrace) String() string {
	return fmt.Sprintf(
		"{seqno:%d, time:%v, dur:%v, error:%v}",
		trace.Step,
		trace.Time,
		trace.Duration,
		trace.Error,
	)
}

/**
	chain消息
 */
type ChainCtx struct {
	context.Context
	mu       sync.Mutex
	done     chan struct{}
	err      error

	Seqno    int64             // 序号
	T        chainMsgType      // 消息类型,main type
	Data     interface{}       // data
	Sync     bool              // 是否等待消息执行结果返回
	AckData  interface{}
							   // TODO 如何更好的返回结果??
	//syncChan chan *ChainMsgACK // 接收结果消息返回channel
	Track    bool         // 是否追踪消息
	Traces   []ChainTrace // 追踪结果

}

func (c *ChainCtx) Done() <-chan struct{} {
	c.mu.Lock()
	if c.done == nil {
		c.done = make(chan struct{})
	}
	d := c.done
	c.mu.Unlock()
	return d
}

func (c *ChainCtx) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

func (c *ChainCtx) Close(ack interface{},err error){
	c.mu.Lock()
	c.AckData = ack
	c.err = err
	if c.done == nil {
		c.done = make(chan struct{})
	}
	close(c.done)
	c.mu.Unlock()
}


func NewContext(t chainMsgType, d interface{}, sync bool, track bool) *ChainCtx {
	msg := &ChainCtx{
		Seqno: time.Now().Unix(),
		T: t,
		Data: d,
		Sync:sync,
		Track:track,
	}
	return msg
}

func NewAddItemContext(d interface{}, sync bool) *ChainCtx {
	return NewContext(CHAIN_ADD_ITEM, d, sync, false)
}

func (c *ChainCtx)String() string {
	return fmt.Sprintf(
		"{seqno:%d, t:%v, sync:%v,track:%v traces:%v}",
		c.Seqno,
		c.T,
		c.Sync,
		c.Track,
		c.Traces,
	)
}
func (c *ChainCtx)SimpleString() string {
	return fmt.Sprintf(
		"{seqno:%d, t:%v, sync:%v,track:%v}",
		c.Seqno,
		c.T,
		c.Sync,
		c.Track,
	)
}