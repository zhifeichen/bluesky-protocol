package models

import (
	"strings"
)

func NewMetricData(metric string, endPoint string, val interface{}, dataType string, tags ...string) *MetaData {
	mv := MetaData{
		Metric:   metric,
		Value:    val,
		Endpoint: endPoint,
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
func GaugeValue(metric, endPoint string, val interface{}, tags ...string) *MetaData {
	return NewMetricData(metric, endPoint, val, METRIC_TYPE_GAUGE, tags...)
}
