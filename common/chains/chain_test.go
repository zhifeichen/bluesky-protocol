package chains

import (
	"fmt"
	"testing"
	"errors"
	"sync"
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

		// 1. 创建任务链
		lineChain := NewLineChains("测试")
		fmt.Println(lineChain.String())
		// 2. 运行任务链
		lineChain.Run()

		// 3. 加入执行任务
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

		// 4. 处理数据,并等待返回

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

		// 4. 处理数据,并等待返回
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

		// 5. 关闭任务链
		lineChain.Stop()

	})
}

func NewPrintIChain(name string, words []string) *LineChain {
	lineChain := NewLineChains(name)
	lineChain.Run()
	for i, w := range words {
		say := NewPrintItem(fmt.Sprintf("%s-%d", name, i), w)
		if err := lineChain.AddTask(say); err != nil {
			fmt.Println("新增item ", say.String(), " err:", err)
		}
	}
	return lineChain
}

func TestTreeChain(t *testing.T) {
	t.Run("测试 tree chain", func(t *testing.T) {
		// 1. 创建根 任务链
		lineChain := NewLineChains("测试root")
		fmt.Println(lineChain.String())

		// 2. 运行根任务链
		lineChain.Run()

		// 3. 创建并加入多个子任务链
		if err := genTreeChains(lineChain); err != nil {
			t.Fail()
		}
		fmt.Println(lineChain.String())

		// 5. 运行数据
		doRunData(lineChain)

		if 0 != lineChain.CountWorks() {
			t.Fatal("任务计数不正确")
		}

		// 5. 关闭任务链
		lineChain.Stop()

	})
}

func genTreeChains(lineChain *LineChain) error {
	say1 := NewPrintIChain("sayIChain-1", []string{" wh", "at"})
	say2 := NewPrintIChain("sayIChain-2", []string{" ", "a"})
	say3 := NewPrintIChain("sayIChain-3", []string{" ni", "ce"})

	if err := lineChain.AddIChain(say1); err != nil {
		fmt.Println("新增iChain ", say1.String(), " err:", err)
		return errors.New("新增任务失败")
	}
	if err := lineChain.AddIChain(say2); err != nil {
		fmt.Println("新增iChain ", say2.String(), " err:", err)
		return errors.New("新增任务失败")
	}
	if err := lineChain.AddIChain(say3); err != nil {
		fmt.Println("新增iChain ", say3.String(), " err:", err)
		return errors.New("新增任务失败")
	}
	// 4. 创建并加入子任务
	say4 := NewPrintItem("sayItem-4", " day!")
	if err := lineChain.AddTask(say4); err != nil {
		fmt.Println("新增iChain ", say4.String(), " err:", err)
		return errors.New("新增任务失败")
	}
	return nil
}

func doRunData(lineChain *LineChain) error {
	data := "i say:"
	if totalWord, traces, err := lineChain.Do(data, true, true); err != nil {
		fmt.Println("处理数据:", data, " 错误:", err)
		return err
	} else {
		fmt.Println("处理数据:\"", data, "\" 结果:\"", totalWord, "\"")
		for _, t := range traces {
			fmt.Print(t.String(), "\n")
		}
		if totalWord != "i say: what a nice day!" {
			return errors.New(fmt.Sprintf("%v != %v", totalWord, "i say: what a nice day! isn't it?"))
		}
		return nil
	}
}

func runData(wg *sync.WaitGroup, lineChain *LineChain) error {
	defer wg.Done()
	return doRunData(lineChain)
}

func TestMultiDatas(t *testing.T) {
	t.Run("测试 tree chain", func(t *testing.T) {
		// 1. 创建根 任务链
		lineChain := NewLineChains("测试root")
		fmt.Println(lineChain.String())

		// 2. 运行根任务链
		lineChain.Run()

		// 3. 创建并加入多个子任务链
		if err := genTreeChains(lineChain); err != nil {
			t.Fail()
		}
		fmt.Println(lineChain.String())

		// 5. 运行数据
		wg := &sync.WaitGroup{}

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go runData(wg, lineChain)
		}

		wg.Wait()
		if 0 != lineChain.CountWorks() {
			t.Fatal("任务计数不正确")
		}

		// 5. 关闭任务链
		lineChain.Stop()

	})
}
