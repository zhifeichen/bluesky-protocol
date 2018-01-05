package jianchi

import (
	"bytes"
	"testing"
)

func TestCRC(t *testing.T) {
	t.Run("request crc", func(t *testing.T) {
		msg := []byte{0x01, 0x03, 0x00, 0x7F, 0x00, 0x05}
		crc := GenCRC(msg)
		if crc != 0x11B4 {
			t.Error("gen crc error: ", crc)
		}
	})
	t.Run("response crc", func(t *testing.T) {
		msg := []byte{0x01, 0x03, 0x0A, 0x00, 0x02, 0x00, 0x00, 0x00, 0x14, 0x00, 0x02, 0x00, 0x08}
		crc := GenCRC(msg)
		if crc != 0xD3AD {
			t.Error("gen crc error: ", crc)
		}
	})
}

func TestRequest(t *testing.T) {
	t.Run("decode", func(t *testing.T) {
		msg := []byte{0x01, 0x03, 0x00, 0x7F, 0x00, 0x05, 0xB4, 0x11}
		r := Request{}
		err := r.Unmarshal(msg)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("encode", func(t *testing.T) {
		r := Request{RequestHeader{1, 3, 0x007F, 0x0005}, 0}
		// r.GenCRC()
		msg, err := r.Marshal()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(msg, []byte{0x01, 0x03, 0x00, 0x7F, 0x00, 0x05, 0xB4, 0x11}) {
			t.Error("encode error: ", msg)
		}
	})
}

func TestResponse(t *testing.T) {
	t.Run("decode", func(t *testing.T) {
		msg := []byte{0x01, 0x03, 0x0A, 0x00, 0x02, 0x00, 0x00,0x00, 0x14, 0x00, 0x02, 0x00, 0x08, 0xAD, 0xD3}
		r := Response{}
		err := r.Unmarshal(msg)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("encode", func(t *testing.T) {
		r := Response{ResponseHeader{1, 3, 0x0A}, []uint16{0x0002, 0x0000, 0x0014, 0x0002, 0x0008}, 0}
		r.GenCRC()
		msg, err := r.Marshal()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(msg, []byte{0x01, 0x03, 0x0A, 0x00, 0x02, 0x00, 0x00,0x00, 0x14, 0x00, 0x02, 0x00, 0x08, 0xAD, 0xD3}) {
			t.Error("encode error: ", msg)
		}
	})
}
