package chains

type ChainMsg struct {
	seqno int64
	t chainMsgType
	d interface{}
	sync bool						// 是否等待消息返回
	syncChan chan ChainMsg
}

func NewMsg(t chainMsgType,d interface{},sync bool) *ChainMsg{
	msg := &ChainMsg{
		t: t,
		d: d,
		sync:sync,
	}
	if sync{
		msg.syncChan = make(chan ChainMsg)
	}
	return msg
}


func NewAddItemMsg(d interface{},sync bool) *ChainMsg{
	return NewMsg(CHAIN_ADD_ITEM,d,sync)
}