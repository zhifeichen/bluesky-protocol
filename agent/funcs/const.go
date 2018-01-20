package funcs

import "encoding/json"

const (
	HTTP_STATE_OK = 200
)

type httpRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"d"`
}

func NewHttpRes(body string) (*httpRes, error) {
	var res httpRes
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return nil, err
	} else {
		return &res, nil
	}

}

func (r *httpRes) CheckRes() bool {
	return r.Code == HTTP_STATE_OK
}
