package bluesky

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type DataUnitHeader struct {
	ID byte
	Count byte
}

type DataUnit24Body struct {
	Body   byte
	UserID byte
	Second uint8
	Minute uint8
	Hour   uint8
	Day    uint8
	Month  uint8
	Year   uint8
}

type DataUnit24 struct {
	Header	DataUnitHeader
	Body	 []DataUnit24Body
}

// TODO: modify unmarshal
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
