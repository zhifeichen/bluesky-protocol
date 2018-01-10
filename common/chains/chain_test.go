package chains

import (
	"fmt"
	"testing"
)

type PrintItem struct {
	BaseTask
	say string
}

func NewPrintItem(name, say string) *PrintItem {
	return &PrintItem{
		BaseTask{
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

		if err := lineChain.AddTask(say1); err != nil {
			fmt.Println("新增item ", say1.String(), " err:", err)
			t.Fail()
		}
		if err := lineChain.AddTask(say2); err != nil {
			fmt.Println("新增item ", say2.String(), " err:", err)
			t.Fail()
		}
		if err := lineChain.AddTask(say3); err != nil {
			fmt.Println("新增item ", say3.String(), " err:", err)
			t.Fail()
		}
		if err := lineChain.AddTask(say4); err != nil {
			fmt.Println("新增item ", say4.String(), " err:", err)
			t.Fail()
		}

		fmt.Println(lineChain.String())

		data := "i say:"

		if totalWord, _, err := lineChain.Do(data, true, false); err != nil {
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

		if err := lineChain.AddTask(say5); err != nil {
			fmt.Println("新增item ", say5.String(), " err:", err)
			t.Fail()
		}
		fmt.Println(lineChain.String())
		if totalWord, traces, err := lineChain.Do(data, true, true); err != nil {
			fmt.Println("处理数据:", data, " 错误:", err)
			t.Fail()
		} else {
			fmt.Println("处理数据:\"", data, "\" 结果:\"", totalWord, "\"")
			for _, t := range traces {
				fmt.Print(t.String(), "\n")
			}
			if totalWord != "i say: what a nice day! isn't it?" {
				t.Fatal(fmt.Sprintf("%v != %v", totalWord, "i say: what a nice day! isn't it?"))
			}
		}

		lineChain.Stop()

	})
}

func NewPrintIChain(name string, words []string, t *testing.T) *LineChain {
	lineChain := NewLineChains(name)
	lineChain.Run()
	for i, w := range words {
		say := NewPrintItem(fmt.Sprintf("%s-%d", name, i), w)
		if err := lineChain.AddTask(say); err != nil {
			fmt.Println("新增item ", say.String(), " err:", err)
			t.Fail()
		}
	}
	return lineChain
}

func TestTreeChain(t *testing.T) {
	t.Run("测试 tree chain", func(t *testing.T) {
		lineChain := NewLineChains("测试root")
		fmt.Println(lineChain.String())
		lineChain.Run()
		say1 := NewPrintIChain("sayIChain-1", []string{" wh", "at"}, t)
		say2 := NewPrintIChain("sayIChain-2", []string{" ", "a"}, t)
		say3 := NewPrintIChain("sayIChain-3", []string{" ni", "ce"}, t)
		say4 := NewPrintIChain("sayIChain-4", []string{" d", "ay!"}, t)

		if err := lineChain.AddIChain(say1); err != nil {
			fmt.Println("新增iChain ", say1.String(), " err:", err)
			t.Fail()
		}
		if err := lineChain.AddIChain(say2); err != nil {
			fmt.Println("新增iChain ", say2.String(), " err:", err)
			t.Fail()
		}
		if err := lineChain.AddIChain(say3); err != nil {
			fmt.Println("新增iChain ", say3.String(), " err:", err)
			t.Fail()
		}
		if err := lineChain.AddIChain(say4); err != nil {
			fmt.Println("新增iChain ", say4.String(), " err:", err)
			t.Fail()
		}

		fmt.Println(lineChain.String())

		data := "i say:"
		if totalWord, traces, err := lineChain.Do(data, true, true); err != nil {
			fmt.Println("处理数据:", data, " 错误:", err)
			t.Fail()
		} else {
			fmt.Println("处理数据:\"", data, "\" 结果:\"", totalWord, "\"")
			for _, t := range traces {
				fmt.Print(t.String(), "\n")
			}
			if totalWord != "i say: what a nice day!" {
				t.Fatal(fmt.Sprintf("%v != %v", totalWord, "i say: what a nice day! isn't it?"))
			}
		}

		lineChain.Stop()

	})
}
