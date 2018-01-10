package chains

/**
	任务处理 interface
 */
type ITask interface {
	GetName() string
	Do(data interface{}) (nextData interface{}, err error)
	String() string
}

/**
	任务链interface, 同时也可以做为子任务链形式处理
 */
type IChain interface {
	GetName() string
	String() string

	AddTask(task ITask) error
	AddIChain(chain IChain) error

	Do(data interface{}, sync, trace bool) (interface{}, []ChainTrace, error)
	Run() error
	Stop() error
}
