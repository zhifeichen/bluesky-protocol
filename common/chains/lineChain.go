package chains

/**
线性运行chain, chain中的任务一个个运行,支持嵌套
*/
import (
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"github.com/zhifeichen/bluesky-protocol/common/utils"
	"sync"
	"time"
	"strings"
	"sync/atomic"
)

type LineChain struct {
	Seqno string         // 序列号
	name  string         // 处理链名称
	items []interface{}  // tasks
	ctxes chan *ChainCtx // 消息队列

	mu   sync.Mutex
	once sync.Once

	// 运行信息
	curWorks int32 // 当前运行的任务数
}

func NewLineChains(name string) *LineChain {
	return &LineChain{
		fmt.Sprintf("lc_%d", time.Now().Unix()),
		name,
		make([]interface{}, 0),
		make(chan *ChainCtx, ITEM_CHANNEL_DEFAULT),
		sync.Mutex{},
		sync.Once{},
		0,
	}
}

func (c *LineChain) GetName() string {
	return c.name
}

func (c *LineChain) CountWorks() int32 {
	return c.curWorks
}

/**
打印整个任务链
*/
func (c *LineChain) String() string {
	names := make([]string, 0)
	for _, c := range c.items {
		switch c := c.(type) {
		case ITask:
			names = append(names, c.String())
		case IChain:
			names = append(names, c.String())
		}
	}
	return fmt.Sprintf(" {lineChains name:%s,tasks:[%s]}", c.name, strings.Join(names, " -> "))
}

/**
增加任务
*/
func (c *LineChain) AddTask(task ITask) error {
	ctx := NewAddItemContext(task, true)
	err, _, _ := c.addHandleCtx(ctx)
	return err
}

func (c *LineChain) AddIChain(chain IChain) error {
	ctx := NewAddItemContext(chain, true)
	err, _, _ := c.addHandleCtx(ctx)
	return err
}

func (c *LineChain) Do(data interface{}, sync, trace bool) (interface{}, []ChainTrace, error) {
	ctx := NewContext(CHAIN_HANDLE_DATA, data, sync, trace)
	err, d, traces := c.addHandleCtx(ctx)
	return d, traces, err
}

/**
启动chain
*/
func (c *LineChain) Run() error {
	c.once.Do(func() {
		xlogger.Info("启动chain:", c.name)
		go c.run()
	})
	return nil
}

func (c *LineChain) Stop() error {
	ctx := NewContext(CHAIN_STOP, nil, true, false)
	err, _, _ := c.addHandleCtx(ctx)
	return err
}

/**
-------------  以下为私有方法  ---------------
*/

func (c *LineChain) addWorksCurWorksCount() {
	atomic.AddInt32(&c.curWorks, 1)
}

func (c *LineChain) decWorksCurWorksCount() {
	atomic.AddInt32(&c.curWorks, -1)
}

func (c *LineChain) addHandleCtx(ctx *ChainCtx) (error, interface{}, []ChainTrace) {
	c.ctxes <- ctx
	if ctx.sync {
		done := ctx.Done()
		<-done
		if err := ctx.Err(); err != nil {
			xlogger.Error("处理消息:", ctx.String(), " 失败:", err)
			return common.NewError(common.CHAIN_HANDLE_MSG_ERROR), nil, nil
		} else {
			return nil, ctx.ackData, ctx.traces
		}
	} else {
		return nil, nil, nil
	}
}

func (c *LineChain) run() {
	for {
		select {
		case ctx := <-c.ctxes:
			switch ctx.t {

			// 增加item
			case CHAIN_ADD_ITEM:
				c.handleAddItemCtx(ctx)

			case CHAIN_HANDLE_DATA:
				c.handleCtx(ctx)

			default:
				if _, stop := c.handleCtlCtx(ctx); stop {
					goto OUT_LOOP
				}
			}
		}
	}

OUT_LOOP:
	xlogger.Warn("退出 chain:", c.name)
	c.once = sync.Once{}
}

/**
增加任务链
拷贝任务链并新增任务, 避免任务链执行过程中变化的问题
*/
func (c *LineChain) handleAddItemCtx(ctx *ChainCtx) error {
	t := make([]interface{}, len(c.items))
	c.mu.Lock()
	copy(t, c.items)
	t = append(t, ctx.data)
	c.items = t
	c.mu.Unlock()
	if ctx.sync {
		ctx.Close(nil, nil)
	}
	return nil
}

/**
处理停止启动等控制消息
*/
func (c *LineChain) handleCtlCtx(ctx *ChainCtx) (err error, stop bool) {
	xlogger.Debug("处理线性chain:", c.Seqno, "指令:", ctx.String())
	stop = false
	switch ctx.t {
	case CHAIN_PAUSE:
		//TODO pause!!
	case CHAIN_STOP:
		stop = true
	default:
		xlogger.Error("处理线性chain:", c.Seqno, "未知指令:", ctx.String())
	}

	ctx.Close(nil, nil)

	return err, stop
}

/**
处理普通消息
*/
func (c *LineChain) handleCtx(ctx *ChainCtx) error {
	go c.doCtx(c.items, ctx)
	return nil
}

/**
处理普通消息
*/
func (c *LineChain) doCtx(items []interface{}, ctx *ChainCtx) error {
	c.addWorksCurWorksCount()
	for i, item := range items {
		switch item := item.(type) {
		case ITask:
			d, trace, err := c.doItem(i, item, ctx.data)
			ctx.traces = append(ctx.traces, trace)
			ctx.duration += trace.Duration
			if err != nil {
				xlogger.Error(ctx.String(), " error:", err)
				ctx.Close(nil, err)
				return err
			} else {
				ctx.data = d
			}
		case IChain:
			d, traces, err := item.Do(ctx.data, true, true)
			var (
				st       int64 = 0
				duration int64 = 0
			)
			for _, t := range traces {
				if 0 == st {
					st = t.Time
				}
				duration += t.Duration
			}

			trace := ChainTrace{
				Name:     item.GetName(),
				Step:     i,
				Time:     st,
				Duration: duration,
				Error:    err,
				Children: traces,
			}
			ctx.traces = append(ctx.traces, trace)
			ctx.duration += trace.Duration
			if err != nil {
				xlogger.Error(ctx.String(), " error:", err)
				ctx.Close(nil, err)
				return err
			} else {
				ctx.data = d
			}
		}

	}
	ctx.Close(ctx.data, nil)
	c.decWorksCurWorksCount()
	return nil
}

func (c *LineChain) doItem(step int, item ITask, data interface{}) (interface{}, ChainTrace, error) {
	var (
		st       int64 = 0
		duration int64 = 0
	)
	st = time.Now().UnixNano() / 1000

	d, err := item.Do(data)

	duration = time.Now().UnixNano()/1000 - st
	trace := ChainTrace{
		Name:     item.GetName(),
		Step:     step,
		Time:     st,
		Duration: duration,
		Error:    err,
	}
	return d, trace, err
}
