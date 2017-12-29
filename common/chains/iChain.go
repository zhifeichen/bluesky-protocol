package chains

type IChain interface {
	GetName() string
	Run() error
	HandleMsg(msg *ChainMsg) error
	String()

}