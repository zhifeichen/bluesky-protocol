package servers

import (
	"net"
	"github.com/zhifeichen/bluesky-protocol/agent/servers/server"
	"github.com/zhifeichen/bluesky-protocol/common/protocol/bluesky"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"errors"
)

type BlueSkyHandler struct {

}

func (h* BlueSkyHandler)Handle(msg interface{}, remote net.Addr) {
	// TODO handler protocol pkt
	comm := msg.(*bluesky.Common)
	if comm.Cmd == 2 && len(comm.Data) == 0 {
		res := &bluesky.Common{
			bluesky.Header{
				comm.MagicNum,
				comm.SerailNo,
				comm.MainVer,
				comm.ClientVer,
				comm.Second,
				comm.Minute,
				comm.Hour,
				comm.Day,
				comm.Month,
				comm.Year,
				[6]byte{},
				[6]byte{},
				0,
				3,
			},
			[]byte{},
			0,
		}
		copy(res.Dst[:], comm.Src[:])
		copy(res.Src[:], comm.Dst[:])
		server.SendMsg(res.GetDst(), res)
	} else {
		xlogger.Errorf("handle msg{%+v} not implement...\n", comm)
	}

}

func (h* BlueSkyHandler)Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return bluesky.Split(data, atEOF)
}

func (h* BlueSkyHandler)Decode(msg []byte) (interface{}, uint64, error) {
	if ok := bluesky.CheckCRC(msg); ok {
		msgComm := bluesky.Common{}
		err := msgComm.Unmarshal(msg)
		if err == nil {
			return &msgComm, msgComm.GetSrc(), nil
		}
		return nil, 0, err
	}
	return nil, 0, errors.New("check CRC failed")
}

func (h* BlueSkyHandler)Encode(msg interface{}) ([]byte, error) {
	comm := msg.(*bluesky.Common)
	return comm.Marshal()
}
