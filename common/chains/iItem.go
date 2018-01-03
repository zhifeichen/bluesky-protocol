package chains

type IItem interface {
	GetName() string
	Do(data interface{}) (err error,nextData interface{})
}
