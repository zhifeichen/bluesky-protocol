package chains

/**
chains 的消息
*/
const (
	CHAIN_ADD_ITEM = iota + 1 // 增加item
	CHAIN_STOP                // 停止chain
	CHAIN_PAUSE               // 暂停chain

	CHAIN_HANDLE_DATA = iota + 100

	ITEM_CHANNEL_DEFAULT = 1024
)

type chainMsgType int32
