package server

import (
	"bufio"
	"github.com/zhifeichen/bluesky-protocol/common/protocol/bluesky"
	"errors"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
)

type BlueSkyCodec struct {
}

func (codec *BlueSkyCodec) GetScanSplitFun() bufio.SplitFunc {
	return bluesky.Split
}
func (codec *BlueSkyCodec) Decode(msg []byte) (interface{}, error) {
	defer func() {
		if p := recover(); p != nil {
			xlogger.Errorf("%s read panics: %v\n", p)
		}
	}()
	if ok := bluesky.CheckCRC(msg);ok{
		msgComm := &bluesky.Common{}
		err := msgComm.Unmarshal(msg)
		return msgComm, err
	} else {
		return nil,errors.New("check CRC failed.")
	}
}

// TODO
// 发送数据,
func (codec *BlueSkyCodec) Encode(msg interface{}) ([]byte, error) {
	return nil, nil
}
