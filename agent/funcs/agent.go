package funcs

import (
	"github.com/zhifeichen/bluesky-protocol/common/models"
	"github.com/zhifeichen/bluesky-protocol/agent/cfg"
)

func AgentState() []*models.MetaData {
	return []*models.MetaData{models.GaugeValue(models.AGENT_STATE_KEY_ALIVE, cfg.Config().EndPoint, 1)}
}
