package bluesky

import (
	"errors"
	"reflect"
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

func (h *DataUnitHeader) Unmarshal(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, errors.New("out of range")
	}
	h.ID = data[0]
	h.Count = data[1]
	return 2, nil
}

func unmarshalBody(data []byte, v interface{}) (advance int, err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		err = errors.New("必须传入指针类型")
		return
	}
	size := rv.Type().Size()
	if len(data) < int(size) {
		err = errors.New("out of range")
		return
	}
	err = binary.Read(bytes.NewBuffer(data), binary.LittleEndian, v)
	advance = int(size)
	return
}

func unmarshalBodies(data []byte, v interface{}) (advance int, err error) {
	t := reflect.ValueOf(v)
	if t.Kind() != reflect.Array {
		err = errors.New("必须传人数组")
	}
	for i := 0; i < t.Len(); i++  {
		advance, err = unmarshalBody(data[advance:], t.Index(i).Addr().Interface())
		if err != nil {
			return
		}
	}
	return
}

func (d *DataUnit24) Unmarshal(data []byte) error {
	advance, err := d.Header.Unmarshal(data)
	if err != nil {
		return err
	}
	d.Body = make([]DataUnit24Body, d.Header.Count)
	advance, err = unmarshalBodies(data[advance:], d.Body)
	return err
}

func (d *DataUnit24) Marshal() ([]byte, error) {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, d.Header)
	for _, b := range d.Body {
		binary.Write(&buffer,binary.LittleEndian, b)
	}
	return buffer.Bytes(), nil
}

func HandleDataUnit24(id byte, unit interface{}) {
	dataUnit := unit.(*DataUnit24)
	fmt.Printf("dataUnit: %v\n", dataUnit)
	// do something with unit
}

func RegisterAll() {
	RegisterUnitHandler(24, &DataUnit24{}, HandleDataUnit24)
}
