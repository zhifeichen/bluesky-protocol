package funcs

import (
	"strings"
	"github.com/zhifeichen/bluesky-protocol/common/models"
)

func NewMetricData(metric string, val interface{}, dataType string, tags ...string) *models.MetaData {
	mv := models.MetaData{
		Metric:   metric,
		Value:    val,
		DataType: dataType,
	}

	size := len(tags)

	if size > 0 {
		mv.Tags = strings.Join(tags, ",")
	}

	return &mv
}

/**
	度量值
 */
func GaugeValue(metric string, val interface{}, tags ...string) *models.MetaData {
	return NewMetricData(metric, val, models.METRIC_TYPE_GAUGE, tags...)
}
