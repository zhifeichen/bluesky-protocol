package bluesky

import (
	"encoding/binary"
	"bytes"
	"unsafe"
)

type Header struct{
	MagicNum  uint16
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
	Dst       [6]byte // 48bit 6byte
	DataLen   uint16 // 应用数据单元长度
	Cmd       uint8  // 0x00: 预留; 0x01: 控制命令; 0x02: 发送数据; 0x03: 确认;0x04: 请求; 0x05: 应答; 0x06: 否认; 0x07~0x7F: 预留; 0x80~0xFF: 用户自定义;
}

type Common struct {
	Header
	data      []byte
	Crc       uint8
}

func (c *Common)Len() int {
	return int(unsafe.Sizeof(c.Header)) + len(c.data) + 3
}

func (c *Common) Unmarshal(data []byte) error {
	//[64 64 0 0 1 1 12 59 12 27 12 17 103 43 0 0 0 0 56 91 1 0 0 0 0 0 2 181 35 35]
	//if len(data) < 30 {
	//	return errors.New("data too short")
	//}
	//i := 2
	//serailNo, err := ReadUint16(data[i:])
	//if err != nil {
	//	return err
	//}
	//c.SerailNo = serailNo
	//i += 2
	//c.MainVer = data[i]
	//i++
	//c.ClientVer = data[i]
	//i++
	//c.Second = data[i]
	//i++
	//c.Minute = data[i]
	//i++
	//c.Hour = data[i]
	//i++
	//c.Day = data[i]
	//i++
	//c.Month = data[i]
	//i++
	//c.Year = data[i]
	//i++
	////src, err := ReadUint48(data[i:])
	////if err != nil {
	////	return err
	////}
	//copy(c.Src[:],data[i:i+7])
	//i += 6
	//copy(c.Dst[:],data[i:i+7])
	//i += 6
	//dataLen, err := ReadUint16(data[i:])
	//if err != nil {
	//	return err
	//}
	//c.DataLen = dataLen
	//i += 2
	//c.cmd = data[i]
	//i++
	//if c.DataLen > 0 {
	//	c.data = make([]byte, c.DataLen)
	//	copy(c.data, data[i: i + int(c.DataLen)])
	//	i += int(c.DataLen)
	//}
	//c.crc = data[i]

	//TODO 注意结构中字节对齐问题
	if err :=binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &c.Header); err != nil{
		return err
	}
	st:=int(unsafe.Sizeof(c.Header))-1
	len:= st +int(c.DataLen)
	if c.DataLen > 0{
		c.data = make([]byte, c.DataLen)

		copy(c.data, data[st: len])
	}
	c.Crc = data[len]
	return nil
}

func (c *Common) Marshal() ([]byte, error) {
	//size := c.Len()
	//buf := make([]byte, size)
	//var i int
	//if int(c.DataLen) != len(c.data) {
	//	return []byte{}, errors.New("invalid Common data")
	//}
	//buf[i] = 0x40
	//i++
	//buf[i] = 0x40
	//i++
	//err := WriteUint16(buf[i:], c.SerailNo)
	//if err != nil {
	//	return []byte{}, err
	//}
	//i += 2
	//buf[i] = c.MainVer
	//i++
	//buf[i] = c.ClientVer
	//i++
	//buf[i] = c.Second
	//i++
	//buf[i] = c.Minute
	//i++
	//buf[i] = c.Hour
	//i++
	//buf[i] = c.Day
	//i++
	//buf[i] = c.Month
	//i++
	//buf[i] = c.Year
	//i++
	////err = WriteUint48(buf[i:], c.Src)
	////if err != nil {
	////	return []byte{}, err
	////}
	//copy(buf[i:],c.Src[:])
	//i += 6
	////err = WriteUint48(buf[i:], c.Dst)
	////if err != nil {
	////	return []byte{}, err
	////}
	//copy(buf[i:],c.Dst[:])
	//i += 6
	//err = WriteUint16(buf[i:], c.DataLen)
	//if err != nil {
	//	return []byte{}, err
	//}
	//i += 2
	//buf[i] = c.Cmd
	//i++
	//if c.DataLen > 0 {
	//	copy(buf[i:], c.data)
	//	i += int(c.DataLen)
	//}
	//buf[i] = c.Crc
	//i++
	//buf[i] = 0x23
	//i++
	//buf[i] = 0x23
	//SetCRC(buf)

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, c.Header)
	binary.Write(&buffer,binary.LittleEndian,c.data)
	binary.Write(&buffer,binary.LittleEndian,GenCrc(buffer.Bytes()[2:]))
	binary.Write(&buffer,binary.LittleEndian,[]byte{0x23,0x23})
	return buffer.Bytes(), nil
}
