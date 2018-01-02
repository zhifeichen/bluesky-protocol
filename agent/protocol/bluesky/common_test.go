package bluesky

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var (
	nilData = "4040000001010c3b0c1b0c11672b00000000385b01000000000002b52323"
	nilC = Common{0, 1, 1, 12, 59, 12, 27, 12, 17, 11111, 88888, 0, 2, []byte{}, 181}
	data = "404003000101103b0c1b0c11672b00000000385b010000000a000218010100003b0c1b0c115f2323"
	dataC = Common{3, 1, 1, 16, 59, 12, 27, 12, 17, 11111, 88888, 10, 2, []byte{24, 1, 1, 0, 0, 59, 12, 27, 12, 17}, 95}
)

func equal(a, b Common) bool {
	return bytes.Equal(a.data, b.data) && a.serailNo == b.serailNo &&
		a.mainVer == b.mainVer && a.clientVer == b.clientVer &&
		a.second == b.second && a.minute == b.minute &&
		a.hour == b.hour && a.day == b.day &&
		a.month == b.month && a.year == b.year &&
		a.src == b.src && a.dst == b.dst &&
		a.dataLen == b.dataLen && a.cmd == b.cmd &&
		a.crc == b.crc
}

func TestUnmarshal(t *testing.T){
	t.Run("data len == 0", func(t *testing.T) {
		binMsg, _ := hex.DecodeString(nilData)
		var c Common
		c.Unmarshal(binMsg)
		if !equal(c, nilC) {
			t.Error(c)
		}
	})
	t.Run("data len != 0", func(t *testing.T) {
		binMsg, _ := hex.DecodeString(data)
		var c Common
		c.Unmarshal(binMsg)
		if !equal(c, dataC) {
			t.Error(c)
		}
	})
}
