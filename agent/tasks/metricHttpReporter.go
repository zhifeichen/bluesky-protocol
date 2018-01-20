package tasks

import (
	"github.com/zhifeichen/bluesky-protocol/common/chains"
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/common/models"
	"errors"
	"github.com/zhifeichen/bluesky-protocol/common/request"
	"encoding/json"
	"bytes"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"github.com/zhifeichen/bluesky-protocol/agent/funcs"
)

const name = "metric.reporter.http"

type HttpReporter struct {
	chains.BaseTask
	url    string
	token  string
	method request.HttpMethod
}

func NewHttpReporter(url, token string, method request.HttpMethod) *HttpReporter {
	return &HttpReporter{
		chains.BaseTask{
			name,
		},
		url,
		token,
		method,
	}
}

func (r *HttpReporter) Do(metric interface{}) (interface{}, error) {
	if meticData, ok := metric.(*models.MetaData); ok {
		if body, err := json.Marshal(meticData); err != nil {
			xlogger.Error(err)
			return nil, err
		} else {
			if res, e := request.Fetch(r.url, request.HTTP_POST,
				bytes.NewReader(body),
				map[string]string{
					"x-access-token": r.token,
				}); e != nil {
				xlogger.Error(e)
				return nil, e
			} else {
				if r, e := funcs.NewHttpRes(res); e != nil {
					return nil, e
				} else {
					if r.CheckRes() {
						return metric, nil
					} else {
						return nil, errors.New(r.Msg)
					}
				}
			}
		}
	} else {
		return nil, errors.New("发送的metric类型错误,必须为*models.MetaData类型")
	}

}

func (r *HttpReporter) String() string {
	return fmt.Sprintf("{task name:%s, %s %s}", r.Name, r.method, r.url)
}
