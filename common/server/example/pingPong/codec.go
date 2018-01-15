package pingPong

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
)

var MAGIC_NUM = []byte("@#")

/**
	数据格式

	2byte		2byte		rest_len bytes
	magic_num	rest_len		body

 */

func Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	hasHeader := false
	for advance = 0; advance < len(data)-1; advance++ {
		if data[advance] == MAGIC_NUM[0] && data[advance+1] == MAGIC_NUM[1] {
			hasHeader = true
			break
		}
	}
	if !hasHeader {
		// 未发现数据头,丢弃数据
		return
	}

	if advance+2+2 > len(data) {
		// 未发现数据长度字段
		return
	}

	var bodyLen int16
	// 读取后续字节长度
	binary.Read(bytes.NewBuffer(data[2:]), binary.LittleEndian, &bodyLen)
	if int(bodyLen+4) > len(data) {
		// 数据不完整
		return
	}
	advance = int(bodyLen) + 4
	token = data[4:advance]
	return
}

type PinPongCodec struct {
}

func (codec *PinPongCodec) GetScanSplitFun() bufio.SplitFunc {
	return Split
}
func (codec *PinPongCodec) Decode(msg []byte) (interface{}, error) {
	body := string(msg)
	xlogger.Debug("接收数据:", body, msg)
	return body, nil
}

// TODO
// 发送数据,
func (codec *PinPongCodec) Encode(msg interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, MAGIC_NUM)
	msgAck := fmt.Sprintf("%s-%v", "pong", msg)
	xlogger.Debug("ack:", msgAck)
	body := []byte(msgAck)
	binary.Write(&buffer, binary.LittleEndian, (int16)(len(body)))
	binary.Write(&buffer, binary.LittleEndian, body)
	return buffer.Bytes(), nil
}
