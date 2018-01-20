package models

import "fmt"

type MetaData struct {
	Metric    string      `json:"metric"`
	Endpoint  string      `json:"endpoint"`
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
	Tags      string      `json:"tags"`
	DataType  string      `json:"dataType"`
}

func (t *MetaData) String() string {
	return fmt.Sprintf("{MetaData Endpoint:%s, Metric:%s,Type:%s, Timestamp:%d, Value:%v, Tags:%v}",
		t.Endpoint, t.Metric, t.DataType, t.Timestamp, t.Value, t.Tags)
}
