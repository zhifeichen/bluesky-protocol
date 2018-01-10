package chains

import "fmt"

/**
 * 所有具体的 item 请继承 baseItem
 */
type BaseTask struct {
	Name string
}

func (c *BaseTask) GetName() string {
	return c.Name
}

/**
 * 具体item需要实现该方法
 */
func (c *BaseTask) Do(data interface{}) (interface{}, error) {
	return data, nil
}

func (c *BaseTask) String() string {
	return fmt.Sprintf("{task name:%s}", c.Name)
}
