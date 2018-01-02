package bluesky

import (
	"errors"
)

type Common struct {
	serailNo  uint16
	mainVer   uint8
	clientVer uint8
	second    uint8
	minute    uint8
	hour      uint8
	day       uint8
	month     uint8
	year      uint8
	src       uint64 // 48bit 6byte
	dst       uint64 // 48bit 6byte
	dataLen   uint16 // 应用数据单元长度
	cmd       uint8  // 0x00: 预留; 0x01: 控制命令; 0x02: 发送数据; 0x03: 确认;0x04: 请求; 0x05: 应答; 0x06: 否认; 0x07~0x7F: 预留; 0x80~0xFF: 用户自定义;
	data      []byte
	crc       uint8
}

func Len(c *Common) int {
	return 30 + len(c.data)
}

func (c *Common) Unmarshal(data []byte) error {
	if len(data) < 30 {
		return errors.New("data too short")
	}
	i := 2
	serailNo, err := ReadUint16(data[i:])
	if err != nil {
		return err
	}
	c.serailNo = serailNo
	i += 2
	c.mainVer = data[i]
	i++
	c.clientVer = data[i]
	i++
	c.second = data[i]
	i++
	c.minute = data[i]
	i++
	c.hour = data[i]
	i++
	c.day = data[i]
	i++
	c.month = data[i]
	i++
	c.year = data[i]
	i++
	src, err := ReadUint48(data[i:])
	if err != nil {
		return err
	}
	c.src = src
	i += 6
	dst, err := ReadUint48(data[i:])
	if err != nil {
		return err
	}
	c.dst = dst
	i += 6
	dataLen, err := ReadUint16(data[i:])
	if err != nil {
		return err
	}
	c.dataLen = dataLen
	i += 2
	c.cmd = data[i]
	i++
	if c.dataLen > 0 {
		c.data = make([]byte, c.dataLen)
		copy(c.data, data[i: i + int(c.dataLen)])
		i += int(c.dataLen)
	}
	c.crc = data[i]
	return nil
}

func (c *Common) Marshal() ([]byte, error) {
	size := Len(c)
	buf := make([]byte, size)
	var i int
	if int(c.dataLen) != len(c.data) {
		return []byte{}, errors.New("invalid Common data")
	}
	buf[i] = 0x40
	i++
	buf[i] = 0x40
	i++
	err := WriteUint16(buf[i:], c.serailNo)
	if err != nil {
		return []byte{}, err
	}
	i += 2
	buf[i] = c.mainVer
	i++
	buf[i] = c.clientVer
	i++
	buf[i] = c.second
	i++
	buf[i] = c.minute
	i++
	buf[i] = c.hour
	i++
	buf[i] = c.day
	i++
	buf[i] = c.month
	i++
	buf[i] = c.year
	i++
	err = WriteUint48(buf[i:], c.src)
	if err != nil {
		return []byte{}, err
	}
	i += 6
	err = WriteUint48(buf[i:], c.dst)
	if err != nil {
		return []byte{}, err
	}
	i += 6
	err = WriteUint16(buf[i:], c.dataLen)
	if err != nil {
		return []byte{}, err
	}
	i += 2
	buf[i] = c.cmd
	i++
	if c.dataLen > 0 {
		copy(buf[i:], c.data)
		i += int(c.dataLen)
	}
	buf[i] = c.crc
	i++
	buf[i] = 0x23
	i++
	buf[i] = 0x23
	return buf, nil
}
