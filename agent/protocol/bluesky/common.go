package bluesky

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
