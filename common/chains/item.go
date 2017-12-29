package chains

/**
 * 所有具体的 item 请继承 baseItem
 */
type BaseItem struct {
	name string
}

func (c *BaseItem) GetName() string{
	return c.name
}

/**
 * 具体item需要实现该方法
 */
func (c *BaseItem) Do(data interface{}) (error,interface{}){
	return nil,data
}

