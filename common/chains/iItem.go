package chains

type IItem interface {
	GetName() string
	Do(data interface{}) (error,nextData interface{})
}
