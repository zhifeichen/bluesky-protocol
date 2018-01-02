package bluesky

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type DataUnit24 struct {
	ID     byte // 24
	Count  byte
	Body   byte
	UserID byte
	Second uint8
	Minute uint8
	Hour   uint8
	Day    uint8
	Month  uint8
	Year   uint8
}

func (d *DataUnit24) Unmarshal(data []byte) error {
	err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, d)
	return err
}

func (d *DataUnit24) Marshal() ([]byte, error) {
	return []byte{}, nil
}

func HandleDataUnit24(id byte, unit interface{}) {
	dataUnit := unit.(*DataUnit24)
	fmt.Printf("dataUnit: %v\n", dataUnit)
	// do something with unit
}

func RegisterAll() {
	RegisterUnitHandler(24, &DataUnit24{}, HandleDataUnit24)
}
