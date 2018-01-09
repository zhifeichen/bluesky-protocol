## Chains

    简单的数据流 流式处理框架
    
### 目录
    
    chains
        |__ ...
        
        
    
### 功能

#### 提供线性流模式处理机制 lineChain

    定义一个 包含一系列数据处理任务 的 任务链, 用于数据流处理, 数据在任务链中将依次处理

```

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

    // 1. 定义处理链
    
    lineChain := NewLineChains("测试")                            
    fmt.Println(lineChain.String())
    
    // 2. 启动处理链
    lineChain.Run()
    
    // 3. 加入数据流处理任务
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
    // 4. 打印链内容
    fmt.Println(lineChain.String())

    // 5. 定义处理数据
    data := "i say:"

    // 6. 处理数据,返回最终处理结果
    if err, totalWord, _ := lineChain.HandleData(data, true, false); err != nil {
        fmt.Println("处理数据:", data, " 错误:", err)
        t.Fail()
    } else {
        fmt.Println("处理数据:\"", data, "\" 结果:\"", totalWord, "\"")
        if totalWord != "i say: what a nice day!" {
            t.Fatal(fmt.Sprintf("%v != %v", totalWord, "i say: what a nice day!"))
        }
    }



```    

    

#### 提供树形流模式处理机制

    定义树状任务处理模式, 数据在任务链中沿树运行
    
    
#### 提供数据运行追踪机制
    
#### 提供任务链池机制
    
#### 提供数据运行实时统计机制    
    
    
    
    