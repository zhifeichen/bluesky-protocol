package jianchi

import (
	"unsafe"
	"errors"
	"bytes"
	"encoding/binary"
)

type RequestHeader struct {
	Addr byte
	Cmd byte
	RegAddr uint16
	RegCnt uint16
}

type Request struct {
	RequestHeader
	CRC uint16
}

type ResponseHeader struct {
	Addr byte
	Cmd byte
	Cnt byte
}

type Response struct {
	ResponseHeader
	Body []uint16
	CRC uint16
}

func (r *Request) GenCRC() error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, r.RequestHeader)
	r.CRC = GenCRC(buf.Bytes())
	return err
}

func (r *Request) Unmarshal(data []byte) error {
	crc := GenCRC(data[:len(data) - 2])
	err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &r.RequestHeader)
	if err != nil {
		return err
	}
	err = binary.Read(bytes.NewBuffer(data[len(data) - 2:]), binary.LittleEndian, &r.CRC)
	if err != nil {
		return err
	}
	if crc != r.CRC {
		return errors.New("crc check error")
	}
	return nil
}

func (r *Request) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, r.RequestHeader)
	err = binary.Write(&buf, binary.LittleEndian, GenCRC(buf.Bytes()))
	return buf.Bytes(), err
}

func (r *Response) GenCRC() error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, r.ResponseHeader)
	if err != nil {
		return err
	}
	if len(r.Body) > 0 {
		err = binary.Write(&buf, binary.LittleEndian, r.Body)
		if err != nil {
			return err
		}
	}
	r.CRC = GenCRC(buf.Bytes())
	return nil
}

func (r *Response) Unmarshal(data []byte) error {
	crc := GenCRC(data[:len(data) - 2])
	err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &r.ResponseHeader)
	if err != nil {
		return err
	}
	if r.Cnt > 0 {
		r.Body = make([]uint16, r.Cnt / 2)
		hLen := int(unsafe.Sizeof(r.ResponseHeader))
		err := binary.Read(bytes.NewBuffer(data[hLen: hLen + int(r.Cnt)]), binary.LittleEndian, r.Body)
		if err != nil {
			return err
		}
	}
	err = binary.Read(bytes.NewBuffer(data[len(data) - 2:]), binary.LittleEndian, &r.CRC)
	if err != nil {
		return err
	}
	if crc != r.CRC {
		return errors.New("crc check error")
	}
	return nil
}

func (r *Response) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, r.ResponseHeader)
	if err != nil {
		return []byte{}, err
	}
	if len(r.Body) > 0 {
		err = binary.Write(&buf, binary.BigEndian, r.Body)
		if err != nil {
			return []byte{}, err
		}
	}
	err = binary.Write(&buf, binary.LittleEndian, GenCRC(buf.Bytes()))
	return buf.Bytes(), err
}

func GenCRC(data []byte) (crc uint16) {
	crc = 0xFFFF
	for _, v := range data {
		crc = crc ^ uint16(v)
		for i := 0; i < 8; i++ {
			if (crc & 0x0001) == 0x0001 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc = crc >> 1
			}
		}
	}
	return
}
