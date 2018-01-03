package chains

/**
	线性运行chain, chain中的任务一个个运行,没有嵌套
 */
import (
	"strings"
	"github.com/zhifeichen/bluesky-protocol/common/logger"
	"github.com/zhifeichen/bluesky-protocol/common/utils"
	"sync"
	"fmt"
	"time"
)

type LineChain struct {
	Seqno string
	name  string
	items []IItem
	msgs  chan *ChainMsg
	once  *sync.Once
}

func NewLineChains(name string) *LineChain {
	return &LineChain{
		fmt.Sprintf("lc_%d", time.Now().Unix()),
		name,
		make([]IItem, 0),
		make(chan *ChainMsg),
		&sync.Once{},
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
func (c *LineChain)AddItems(item  IItem) error {
	msg := NewAddItemMsg(item, false)
	err, _, _ := c.addHandleMsg(msg)
	return err
}
func (c *LineChain)AddSyncItem(item IItem) error {
	msg := NewAddItemMsg(item, true)
	err, _, _ := c.addHandleMsg(msg)
	return err
}

func (c *LineChain)HandleData(data interface{},sync,trace bool) (error,interface{},[]ChainMsgTrace){
	msg := NewMsg(CHAIN_HANDLE_DATA, data, sync, trace)
	err, d ,traces:= c.addHandleMsg(msg)
	return err,d, traces
}

func (c *LineChain) addHandleMsg(msg *ChainMsg) (error, interface{},[]ChainMsgTrace) {
	c.msgs <- msg
	if msg.Sync && msg.syncChan != nil {
		if msgAck, ok := <-msg.syncChan; !ok {
			logger.Info.Println("处理消息:", msg.String()," 失败")
			return common.NewError(common.CHAIN_HANDLE_MSG_ERROR), nil,nil
		} else {
			logger.Info.Println("处理消息: seqno:", msg.Seqno," 成功 ",msgAck)
			return nil, msgAck.Data,msg.Traces
		}

	} else {
		return nil, nil, nil
	}
}

/**
	启动chain
 */
func (c *LineChain) Run() {
	c.once.Do(func() {
		logger.Info.Println("启动chain:", c.name)
		go c.run()
	})

}

func (c *LineChain)run() {
	for {
		select {
		case msg := <-c.msgs:
			switch msg.T {

			// 增加item
			case CHAIN_ADD_ITEM:
				c.handleAddItemMsg(msg)

			case CHAIN_HANDLE_DATA:
				c.handleMsg(msg)

			default:
				if _, stop := c.handleCtlMsg(msg); stop {
					goto OUT_LOOP
				}
			}
		}
	}

	OUT_LOOP:
	logger.Warning.Println("退出循环chain:", c.name)
	c.once = &sync.Once{}
}


/**
	增加任务链
	拷贝任务链并新增任务, 避免任务链执行过程中变化的问题
 */
func (c *LineChain)handleAddItemMsg(msg *ChainMsg) error {
	t := make([]IItem, len(c.items))
	copy(t, c.items)
	t = append(t, msg.Data.(IItem))
	c.items = t
	if msg.Sync && msg.syncChan != nil {
		msg.syncChan <- NewMsgAck(msg.Seqno, msg.T, nil, nil)
	}
	return nil
}

/**
	处理停止启动等控制消息
 */
func (c *LineChain)handleCtlMsg(msg *ChainMsg) (err error, stop bool) {
	logger.Info.Println("处理线性chain:", c.Seqno, "指令:", msg.String())
	stop = false
	switch msg.T {
	case CHAIN_PAUSE:
	//TODO pause!!
	case CHAIN_STOP:
		stop = true
	default:
		logger.Error.Println("处理线性chain:", c.Seqno, "未知指令:", msg.String())
	}

	if msg.Sync {
		msg.syncChan <- NewMsgAck(msg.Seqno, msg.T, nil, nil)
	}

	return err, stop
}

/**
	处理普通消息
 */
func (c *LineChain)handleMsg(msg *ChainMsg) error {
	go c.doMsg(c.items, msg)
	return nil
}
/**
	处理普通消息
 */
func (c *LineChain)doMsg(items []IItem, msg *ChainMsg) error {
	for i, item := range items {
		var st int64 = 0
		if msg.Track {
			st = time.Now().UnixNano() / 1000
		}

		err, d := item.Do(msg.Data)
		if msg.Track{
			trace := ChainMsgTrace{
				i,
				st,
				time.Now().UnixNano() / 1000 - st,
				err,
			}
			msg.Traces = append(msg.Traces, trace)
		}
		if err != nil {
			logger.Error.Println(msg.String(), " error:", err)
			if msg.Sync && msg.syncChan != nil {
				msg.syncChan <- NewMsgAck(msg.Seqno, msg.T, nil, err)
			}
			return err
		} else {
			msg.Data = d
		}
	}
	if msg.Sync && msg.syncChan != nil {
		msg.syncChan <- NewMsgAck(msg.Seqno, msg.T, msg.Data, nil)
	}

	return nil
}
