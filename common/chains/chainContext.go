package chains

import (
	"context"
	"fmt"
	"sync"
	"time"
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
	mu   sync.Mutex
	done chan struct{}
	err  error

	seqno   int64        // 序号
	t       chainMsgType // 消息类型,main type
	data    interface{}  // data
	ackData interface{}

	sync   bool         // 是否等待消息执行结果返回
	track  bool         // 是否追踪消息
	traces []ChainTrace // 追踪结果

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

func (c *ChainCtx) Close(ack interface{}, err error) {
	c.mu.Lock()
	c.ackData = ack
	c.err = err
	if c.done == nil {
		c.done = make(chan struct{})
	}
	close(c.done)
	c.mu.Unlock()
}

func NewContext(t chainMsgType, d interface{}, sync bool, track bool) *ChainCtx {
	msg := &ChainCtx{
		seqno: time.Now().Unix(),
		t:     t,
		data:  d,
		sync:  sync,
		track: track,
	}
	return msg
}

func NewAddItemContext(d interface{}, sync bool) *ChainCtx {
	return NewContext(CHAIN_ADD_ITEM, d, sync, false)
}

func (c *ChainCtx) String() string {
	return fmt.Sprintf(
		"{seqno:%d, t:%v, sync:%v,track:%v traces:%v}",
		c.seqno,
		c.t,
		c.sync,
		c.track,
		c.traces,
	)
}
func (c *ChainCtx) SimpleString() string {
	return fmt.Sprintf(
		"{seqno:%d, t:%v, sync:%v,track:%v}",
		c.seqno,
		c.t,
		c.sync,
		c.track,
	)
}
