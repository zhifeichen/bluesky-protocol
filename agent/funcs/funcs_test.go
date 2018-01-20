package funcs

import (
	"testing"
	"fmt"
	"net/url"
	"encoding/json"
)

func TestChain(t *testing.T) {
	t.Run("测试 func", func(t *testing.T) {
		var m map[string]string
		fmt.Println(m == nil, len(m))
		for k, v := range m {
			fmt.Println(k, v)
		}

		fmt.Println(url.QueryEscape("t=百度"))
		pmap := map[string]string{
			"t": "百度",
		}
		b, _ := json.Marshal(pmap)
		fmt.Println(string(b))
	})
}
