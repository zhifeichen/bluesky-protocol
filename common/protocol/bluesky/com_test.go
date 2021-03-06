package bluesky

import (
	"testing"
	"encoding/hex"
	"encoding/binary"
	"bytes"
	"fmt"
	"unsafe"
)

type Com struct {
	MargicNum uint16
	SerailNo  uint16
	MainVer   uint8
	ClientVer uint8
	Second    uint8
	Minute    uint8
	Hour      uint8
	Day       uint8
	Month     uint8
	Year      uint8
	Src       [6]byte // 48bit 6byte
	Dst       [6]byte  // 48bit 6byte
	DataLen   uint16 // 应用数据单元长度
	Cmd       uint8
	//Data      []byte
	//Crc       uint8
}

func TestDecode(t *testing.T){
	msg := "4040000001010c3b0c1b0c11672b00000000385b01000000000002b52323"
	//msg := "404003000101103b0c1b0c11672b00000000385b010000000a000218010100003b0c1b0c115f2323"
	data, _ := hex.DecodeString(msg)
	t.Run("data len != 0", func(t *testing.T) {
		cm := Com{}
		fmt.Println(data,int(unsafe.Sizeof(cm)+12))
		err :=binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &cm)
		fmt.Println(err,"cm:",cm)
	})
}

func TestType(t *testing.T) {
	RegisterAll()
	err := HandleRawData([]byte{24, 1, 1, 0, 0, 59, 12, 27, 12, 17})
	if err != nil {
		t.Error(err)
	}
	d := DataUnit24{}
	err = d.Unmarshal([]byte{24, 1, 1, 0, 0, 59, 12, 27, 12, 17})
	if err != nil {
		t.Error(d)
	}
	r, err := d.Marshal()
	if !bytes.Equal(r, []byte{24, 1, 1, 0, 0, 59, 12, 27, 12, 17}) {
		t.Error(r)
	}
}
