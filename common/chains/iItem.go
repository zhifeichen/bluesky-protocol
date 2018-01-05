package chains

type IItem interface {
	GetName() string
	Do(data interface{}) (nextData interface{},err error)
	String() string
}
