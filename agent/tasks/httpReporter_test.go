package tasks

import (
	"testing"
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/common/chains"
	"github.com/zhifeichen/bluesky-protocol/common/request"
	"github.com/zhifeichen/bluesky-protocol/common/models"
	"time"
	"github.com/zhifeichen/bluesky-protocol/common/utils"
)

const (
	uploadUrl = "http://127.0.0.1:3004/agent/uploadServices/uploadMetric"
	token     = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhY2NvdW50SWQiOiI1YTYxNTYwOTU5ODJlNzQ5YTA3MjVmNGEiLCJ1c2VybmFtZSI6Imh6amd6MCIsInJvbGUiOjgwLCJzaW5nbGUiOjAsImV4cCI6MTUxNzE5NDIxMDIzMywiaWF0IjoxNTE2MzMwMjEwfQ.TfC1a1efOaQjC5HJWNcGM8nGw_yzKUq_mBtqPRHyfuw"
)

func TestChain(t *testing.T) {
	t.Run("测试 http reporter", func(t *testing.T) {
		lineChain := chains.NewLineChains("测试")

		// 2. 运行任务链
		lineChain.Run()
		defer lineChain.Stop()

		httpReporter := NewHttpReporter(uploadUrl, token, request.HTTP_POST)

		if err := lineChain.AddTask(httpReporter); err != nil {
			fmt.Println("新增item ", httpReporter.String(), " err:", err)
			t.Fail()
		}
		fmt.Println(lineChain.String())
		tags := map[string]string{
			"cpu":   "0",
			"level": "1",
			"AA":    "web",
		}

		m := models.MetaData{
			Metric:    "temperature",
			Endpoint:  "xxx设备",
			Timestamp: time.Now().Unix() * 1000, // ms
			Value:     1,
			Tags:      utils.SortedTags(tags),
			DataType:  models.METRIC_TYPE_GAUGE,
		}

		if res, traces, err := lineChain.Do(&m, true, true); err != nil {
			fmt.Println("处理数据:", m.String(), " 错误:", err)
			t.Fail()
		} else {
			fmt.Println("返回:", res, traces)
		}

	})
}
