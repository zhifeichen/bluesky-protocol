package chains

import "fmt"

/**
 * 所有具体的 item 请继承 baseItem
 */
type BaseItem struct {
	Name string
}

func (c *BaseItem) GetName() string{
	return c.Name
}

/**
 * 具体item需要实现该方法
 */
func (c *BaseItem) Do(data interface{}) (error,interface{}){
	return nil,data
}

func (c *BaseItem) String() string{
	return fmt.Sprintf("{name:%s}",c.Name)
}