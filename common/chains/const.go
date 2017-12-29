package chains

/**
	chains 的消息
 */
const (
	CHAIN_ADD_ITEM			= iota+1				// 增加item
	CHAIN_STOP										// 停止chain
	CHAIN_PAUSE										// 暂停chain
	CHAIN_RUN										// 开始chain
)

type chainMsgType int32