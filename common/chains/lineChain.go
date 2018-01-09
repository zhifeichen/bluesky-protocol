package chains

/**
线性运行chain, chain中的任务一个个运行,没有嵌套
*/
import (
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"github.com/zhifeichen/bluesky-protocol/common/utils"
	"strings"
	"sync"
	"time"
)

type LineChain struct {
	Seqno string         // 序列号
	name  string         // 处理链名称
	items []IItem        // tasks
	ctxes chan *ChainCtx // 消息队列

	mu   sync.Mutex
	once sync.Once
}

func NewLineChains(name string) *LineChain {
	return &LineChain{
		fmt.Sprintf("lc_%d", time.Now().Unix()),
		name,
		make([]IItem, 0),
		make(chan *ChainCtx, ITEM_CHANNEL_DEFAULT),
		sync.Mutex{},
		sync.Once{},
	}
}

/**
打印整个任务链
*/
func (c *LineChain) String() string {
	names := make([]string, 0)
	for _, v := range c.items {
		names = append(names, v.GetName())
	}
	return fmt.Sprintf(" 线性处理链{Name:%s,tasks:[%s]}", c.name, strings.Join(names, " -> "))
}

/**
增加任务
TODO 增加多个任务?
*/
func (c *LineChain) AddItems(item IItem) error {
	ctx := NewAddItemContext(item, false)
	err, _, _ := c.addHandleCtx(ctx)
	return err
}
func (c *LineChain) AddSyncItem(item IItem) error {
	ctx := NewAddItemContext(item, true)
	err, _, _ := c.addHandleCtx(ctx)
	return err
}

func (c *LineChain) HandleData(data interface{}, sync, trace bool) (error, interface{}, []ChainTrace) {
	ctx := NewContext(CHAIN_HANDLE_DATA, data, sync, trace)
	err, d, traces := c.addHandleCtx(ctx)
	return err, d, traces
}

/**
启动chain
*/
func (c *LineChain) Run() {
	c.once.Do(func() {
		xlogger.Info("启动chain:", c.name)
		go c.run()
	})

}

func (c *LineChain) Stop() error {
	ctx := NewContext(CHAIN_STOP, nil, true, false)
	err, _, _ := c.addHandleCtx(ctx)
	return err
}

/**
-------------  以下为私有方法  ---------------
*/

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
	t := make([]IItem, len(c.items))
	c.mu.Lock()
	copy(t, c.items)
	t = append(t, ctx.data.(IItem))
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
func (c *LineChain) doCtx(items []IItem, ctx *ChainCtx) error {
	for i, item := range items {
		var (
			st       int64 = 0
			duration int64 = 0
		)
		st = time.Now().UnixNano() / 1000

		d, err := item.Do(ctx.data)

		duration = time.Now().UnixNano()/1000 - st
		trace := ChainTrace{
			i,
			st,
			duration,
			err,
		}
		ctx.traces = append(ctx.traces, trace)
		ctx.duration += duration
		if err != nil {
			xlogger.Error(ctx.String(), " error:", err)
			ctx.Close(nil, err)
			return err
		} else {
			ctx.data = d
		}
	}
	ctx.Close(ctx.data, nil)

	return nil
}
