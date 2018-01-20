package funcs

import (
	"github.com/zhifeichen/bluesky-protocol/common/models"
	"github.com/zhifeichen/bluesky-protocol/agent/cfg"
)

type MetricFunc func() []*models.MetaData

type MetricProcess struct {
	Fs       []MetricFunc // 执行的方法列表
	Interval int          // 间隔时间
}

func GenMetircProcesses() []MetricProcess {
	interval := cfg.Config().Transfer.Interval
	return []MetricProcess{
		{
			[]MetricFunc{
				AgentState,
			},
			interval,
		},
	}
}
