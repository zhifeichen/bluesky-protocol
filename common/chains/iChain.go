package chains

type IChain interface {
	GetName() string
	Run() error
	HandleMsg(msg *ChainCtx) error
	String()
}
