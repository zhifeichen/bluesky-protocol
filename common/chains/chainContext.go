package chains

import (
	"fmt"
	"time"
	"github.com/zhifeichen/bluesky-protocol/common/context"
)

/**
追踪消息
*/
type ChainTrace struct {
	Name     string
	Step     int
	Time     int64 // ms
	Duration int64
	Error    error
	Children []ChainTrace
}

func (trace *ChainTrace) String() string {

	msgs := make([]string, len(trace.Children)+2)
	for _, t := range trace.Children {
		msgs = append(msgs, t.String())
	}
	return fmt.Sprintf(
		"{seqno:%d,name:%s, time:%v, dur:%vus, err:%v, subs:%v}",
		trace.Step,
		trace.Name,
		trace.Time,
		trace.Duration,
		trace.Error,
		msgs,
	)

}

/**
chain消息
*/
type ChainCtx struct {
	context.Context
	ctxCancel context.CancelFunc
	seqno     int64        // 序号
	t         chainMsgType // 消息类型,main type
	data      interface{}  // data
	ackData   interface{}

	sync     bool         // 是否等待消息执行结果返回
	track    bool         // 是否追踪消息
	traces   []ChainTrace // 追踪结果
	duration int64
}

func (c *ChainCtx) Close(ack interface{}, err error) {
	c.Context.Close(err)
	c.ackData = ack
}

func NewContext(t chainMsgType, d interface{}, sync bool, track bool) *ChainCtx {
	ctx, ctxCancel := context.WithCancel(context.Background())
	msg := &ChainCtx{
		Context:   ctx,
		ctxCancel: ctxCancel,
		seqno:     time.Now().Unix(),
		t:         t,
		data:      d,
		sync:      sync,
		track:     track,
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
