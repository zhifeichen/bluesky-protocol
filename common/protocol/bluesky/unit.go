package bluesky

import (
	"fmt"
	"errors"
	"reflect"
)

type DataUnitHandler func(typeId byte, unit interface{})

type DataUnitInterface interface {
	Unmarshal(data []byte) error
	Marshal() ([]byte, error)
}

type DataUnit struct {
	dataType	reflect.Type
	handler		DataUnitHandler
}

var (
	unitMap = make(map[byte]DataUnit)
)

func RegisterUnitHandler(dataTypeId byte, inter interface{}, handler DataUnitHandler) {
	var unit DataUnit
	unit.dataType = reflect.ValueOf(inter).Elem().Type()
	unit.handler = handler

	fmt.Println("data type: ", unit.dataType.String())
	dataUnitType := reflect.TypeOf((*DataUnitInterface)(nil)).Elem()
	fmt.Println("data unit type: ", dataUnitType.String())
	// if !unit.dataType.Implements(dataUnitType) {
	// 	panic("DataUnit should implements method marshal and unmarshal.")
	// }

	unitMap[dataTypeId] = unit
}

func handleRawData(dataTypeId byte, data []byte) error {
	if unit, ok := unitMap[dataTypeId]; ok {
		msg := reflect.New(unit.dataType).Interface().(DataUnitInterface)
		err := msg.Unmarshal(data)
		if err != nil {
			return err
		}
		unit.handler(dataTypeId, msg)
		return err
	}
	return errors.New("not found dataTypeId")
}

func HandleRawData(data []byte) error {
	return handleRawData(data[0], data)
}
