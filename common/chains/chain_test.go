package chains

import (
	"fmt"
	"testing"
)

type PrintItem struct {
	BaseItem
	say string
}

func NewPrintItem(name, say string) *PrintItem {
	return &PrintItem{
		BaseItem{
			Name: name,
		},
		say,
	}
}

func (item *PrintItem) Do(data interface{}) (interface{}, error) {
	words := data.(string) + item.say
	fmt.Println(item.GetName(), " say:", words)
	return words, nil
}

func TestChain(t *testing.T) {
	t.Run("测试 chain", func(t *testing.T) {
		lineChain := NewLineChains("测试")
		fmt.Println(lineChain.String())
		lineChain.Run()
		say1 := NewPrintItem("sayItem-1", " what")
		say2 := NewPrintItem("sayItem-2", " a")
		say3 := NewPrintItem("sayItem-3", " nice")
		say4 := NewPrintItem("sayItem-4", " day!")

		if err := lineChain.AddSyncItem(say1); err != nil {
			fmt.Println("新增item ", say1.String(), " err:", err)
			t.Fail()
		}
		if err := lineChain.AddSyncItem(say2); err != nil {
			fmt.Println("新增item ", say1.String(), " err:", err)
			t.Fail()
		}
		if err := lineChain.AddSyncItem(say3); err != nil {
			fmt.Println("新增item ", say1.String(), " err:", err)
			t.Fail()
		}
		if err := lineChain.AddSyncItem(say4); err != nil {
			fmt.Println("新增item ", say1.String(), " err:", err)
			t.Fail()
		}

		fmt.Println(lineChain.String())

		data := "i say:"

		if err, totalWord, _ := lineChain.HandleData(data, true, false); err != nil {
			fmt.Println("处理数据:", data, " 错误:", err)
			t.Fail()
		} else {
			fmt.Println("处理数据:\"", data, "\" 结果:\"", totalWord, "\"")
			if totalWord != "i say: what a nice day!" {
				t.Fatal(fmt.Sprintf("%v != %v", totalWord, "i say: what a nice day!"))
			}
		}

		fmt.Println("")

		say5 := NewPrintItem("sayItem-5", " isn't it?")

		if err := lineChain.AddSyncItem(say5); err != nil {
			fmt.Println("新增item ", say5.String(), " err:", err)
			t.Fail()
		}
		fmt.Println(lineChain.String())
		if err, totalWord, traces := lineChain.HandleData(data, true, true); err != nil {
			fmt.Println("处理数据:", data, " 错误:", err)
			t.Fail()
		} else {
			fmt.Println("处理数据:\"", data, "\" 结果:\"", totalWord, "\" tracks:", traces)
			if totalWord != "i say: what a nice day! isn't it?" {
				t.Fatal(fmt.Sprintf("%v != %v", totalWord, "i say: what a nice day! isn't it?"))
			}
		}

		lineChain.Stop()

	})
}
