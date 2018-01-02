package bluesky_test

import (
	"fmt"
	"testing"
	"github.com/zhifeichen/bluesky-protocol/protocol/bluesky"
)

var (
	testData []byte = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	testBuf = make([]byte, 8)
)

func TestRead(t *testing.T) {
	t.Run("read uint16", func(t *testing.T) {
		ret, err := bluesky.ReadUint16(testData)
		if err != nil {
			t.Error(err.Error())
		}
		if ret != 0x0201 {
			t.Errorf("except 0x0201,but read: %d", ret)
		}
	})

	t.Run("read uint32", func(t *testing.T) {
		ret, err := bluesky.ReadUint32(testData)
		if err != nil {
			t.Error(err.Error())
		}
		if ret != 0x04030201 {
			t.Errorf("except 0x04030201,but read: %d", ret)
		}
	})

	t.Run("read uint48", func(t *testing.T) {
		ret, err := bluesky.ReadUint48(testData)
		if err != nil {
			t.Error(err.Error())
		}
		if ret != 0x060504030201 {
			t.Errorf("except 0x060504030201,but read: %d", ret)
		}
	})

	t.Run("read uint64", func(t *testing.T) {
		ret, err := bluesky.ReadUint64(testData)
		if err != nil {
			t.Error(err.Error())
		}
		if ret != 0x0807060504030201 {
			t.Errorf("except 0x0807060504030201,but read: %d", ret)
		}
	})
}

func checkBuf(t *testing.T, buf []byte, cap int) {
	for i := 0; i < cap; i++ {
		if buf[i] != byte(i) + 1 {
			t.Errorf("except buf[%d]: %d, but get %d", i, i + 1, buf[i])
		}
	}
}

func TestWrite(t *testing.T) {
	t.Run("write uint16", func(t *testing.T) {
		testBuf = make([]byte, 8)
		err := bluesky.WriteUint16(testBuf, 0x0201)
		if err != nil {
			t.Error(err.Error())
		}
		fmt.Printf("testBuf: %v\n", testBuf)
		checkBuf(t, testBuf, 2)
	})

	t.Run("write uint32", func(t *testing.T) {
		testBuf = make([]byte, 8)
		err := bluesky.WriteUint32(testBuf, 0x04030201)
		if err != nil {
			t.Error(err.Error())
		}
		fmt.Printf("testBuf: %v\n", testBuf)
		checkBuf(t, testBuf, 4)
	})

	t.Run("write uint48", func(t *testing.T) {
		testBuf = make([]byte, 8)
		err := bluesky.WriteUint48(testBuf, 0x060504030201)
		if err != nil {
			t.Error(err.Error())
		}
		fmt.Printf("testBuf: %v\n", testBuf)
		checkBuf(t, testBuf, 6)
	})

	t.Run("write uint64", func(t *testing.T) {
		testBuf = make([]byte, 8)
		err := bluesky.WriteUint64(testBuf, 0x0807060504030201)
		if err != nil {
			t.Error(err.Error())
		}
		fmt.Printf("testBuf: %v\n", testBuf)
		checkBuf(t, testBuf, 8)
	})
}
