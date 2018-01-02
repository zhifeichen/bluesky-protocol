package bluesky

import (
	"fmt"
	"bufio"
	"errors"
	"io"
)

const (
	START_SYMBLE byte = 0x40
	END_SYMBLE byte = 0x23
)

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for advance = 0; advance < len(data); advance++ {
		if data[advance] == START_SYMBLE && data[advance + 1] == START_SYMBLE {
			break
		}
	}
	if advance == len(data) {
		err = errors.New("out of range")
		return
	}
	for restart := true; restart; {
		start = advance
		for advance += 2; advance < len(data); advance++ {
			if data[advance] == END_SYMBLE && data[advance + 1] == END_SYMBLE {
				restart = false
				break
			}
			if data[advance] == START_SYMBLE && data[advance + 1] == START_SYMBLE {
				break
			}
		}
		if advance >= len(data) - 1 {
			err = errors.New("out of range")
			return
		}
	}
	advance += 2
	token = data[start: advance]
	return
}

func Get(r io.Reader, msgChan chan []byte) {
	scanner := bufio.NewScanner(r)
	scanner.Split(split)
	for scanner.Scan() {
		msg:= scanner.Bytes()
		fmt.Println("scan: ",msg)
		msgChan <- scanner.Bytes()
	}
	msgChan <- []byte{}
}

func CheckCRC(msg []byte) bool {
	var sum byte
	msgLen := len(msg)
	for i := 2; i < msgLen - 3; i++ {
		sum += msg[i]
	}
	return sum == msg[msgLen - 3]
}

func readUintN(data []byte, cap int) (interface{}, error) {
	if len(data) < cap {
		return 0, errors.New("out of range")
	}
	ret := uint64(data[cap - 1])
	for i := cap - 2; i >= 0; i-- {
		ret = ret << 8 + uint64(data[i])
	}
	return ret, nil
}

func ReadUint16(data []byte) (uint16, error) {
	// if len(data) < 2 {
	// 	return 0, errors.New("out of range")
	// }
	// ret := uint16(data[1])
	// ret = ret << 8 + uint16(data[0])
	ret, err := readUintN(data, 2)
	return uint16(ret.(uint64) & 0xFFFF), err
}

func ReadUint32(data []byte) (uint32, error) {
	ret, err := readUintN(data, 4)
	return uint32(ret.(uint64) & 0xFFFFFFFF), err
}

func ReadUint48(data []byte) (uint64, error) {
	ret, err := readUintN(data, 6)
	return ret.(uint64), err
}

func ReadUint64(data []byte) (uint64, error) {
	ret, err := readUintN(data, 8)
	return ret.(uint64), err
}

func writeUintN(data []byte, value interface{}, cap int) error {
	if len(data) < cap {
		return errors.New("out of range")
	}
	var ret uint64
	switch value.(type) {
	case uint16:
		ret = uint64(value.(uint16))
		break
	case uint32:
		ret = uint64(value.(uint32))
		break
	case uint64:
		ret = value.(uint64)
		break
	}
	fmt.Printf("value: %d;cap: %d; ", ret, cap)
	for i := uint(0); i < uint(cap); i++ {
		v := byte(ret & 0xFF)
		data[i] = v
		ret = ret >> 8
	}
	fmt.Printf("data: %v\n", data)
	return nil
}

func WriteUint16(data []byte, value uint16) error {
	return writeUintN(data, value, 2)
}

func WriteUint32(data []byte, value uint32) error {
	return writeUintN(data, value, 4)
}

func WriteUint48(data []byte, value uint64) error {
	return writeUintN(data, value, 6)
}

func WriteUint64(data []byte, value uint64) error {
	return writeUintN(data, value, 8)
}
