package bluesky

import (
	"bytes"
	"encoding/hex"
	"testing"
	"fmt"
)

var (
	nilData = "4040000001010c3b0c1b0c11672b00000000385b01000000000002b52323"
	nilC = Common{Header{16448,0, 1, 1, 12, 59, 12, 27, 12, 17, [6]byte{103,43,0,0,0,0}, [6]byte{56,91,1,0,0,0}, 0, 2}, []byte{}, 181}
	data = "404003000101103b0c1b0c11672b00000000385b010000000a000218010100003b0c1b0c115f2323"
	dataC = Common{Header{16448,3, 1, 1, 16, 59, 12, 27, 12, 17, [6]byte{103,43,0,0,0,0}, [6]byte{56,91,1,0,0,0}, 10, 2}, []byte{24, 1, 1, 0, 0, 59, 12, 27, 12, 17}, 95}
)

func equal(a, b Common) bool {
	return bytes.Equal(a.data, b.data) && a.SerailNo == b.SerailNo &&
		a.MainVer == b.MainVer && a.ClientVer == b.ClientVer &&
		a.Second == b.Second && a.Minute == b.Minute &&
		a.Hour == b.Hour && a.Day == b.Day &&
		a.Month == b.Month && a.Year == b.Year &&
		a.Src == b.Src && a.Dst == b.Dst &&
		a.DataLen == b.DataLen && a.Cmd == b.Cmd &&
		a.Crc == b.Crc
}

func TestUnmarshal(t *testing.T){
	t.Run("data len == 0", func(t *testing.T) {
		binMsg, _ := hex.DecodeString(nilData)

		fmt.Println(binMsg)
		var c Common
		c.Unmarshal(binMsg)

		if !equal(c, nilC) {
			t.Error(c)
		}
	})
	t.Run("data len != 0", func(t *testing.T) {
		binMsg, _ := hex.DecodeString(data)
		var c Common
		fmt.Println(binMsg)
		c.Unmarshal(binMsg)
		fmt.Println(c,"\n",dataC)
		if !equal(c, dataC) {
			t.Error(c)
		}
	})
}

func TestMarshal(t *testing.T) {
	t.Run("data len == 0", func(t *testing.T) {
		msg, err := nilC.Marshal()
		if err != nil {
			t.Error(err)
		}
		binMsg, _ := hex.DecodeString(nilData)
		if !bytes.Equal(msg, binMsg) {
			t.Error(err.Error())
		}
	})
	t.Run("data len != 0", func(t *testing.T) {
		msg, err := dataC.Marshal()
		if err != nil {
			t.Error(err)
		}
		binMsg, err := hex.DecodeString(data)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(msg, binMsg) {
			t.Error(msg, binMsg)
		}
	})
}


