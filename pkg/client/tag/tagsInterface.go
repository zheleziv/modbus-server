package tag

import (
	myerr "zheleznovux.com/modbus-console/pkg"
)

// опустить абстракцию ниже в value
type TagInterface interface { // виртуальные поля, которыми должны обладать все тэги
	SetName(string) error // имя
	Name() string

	SetAddress(uint16) error // адрес
	Address() uint16

	SetScanPeriod(float64) error // период сканирования
	ScanPeriod() float64

	DataType() string // строка указывающая на тип данных тэга

	SetState(bool) // состояние подключения
	State() bool

	SetTimestamp() // временная метка
	Timestamp() string

	ReadDevice() error
}

func NewTag(name string, address uint16, scanPeriod float64, dataType string) (TagInterface, error) {
	var tagI TagInterface
	var err error

	if dataType == "" {
		return nil, myerr.New("invalid tag dataType")
	}
	switch dataType {
	case COIL_TYPE:
		tagI, err = NewCoilTagWithData(name, address, scanPeriod)
	case WORD_TYPE:
		tagI, err = NewWordTagWithData(name, address, scanPeriod)
	case DWORD_TYPE:
		tagI, err = NewDWordTagWithData(name, address, scanPeriod)
	default:
		return nil, myerr.New("invalid tag dataType")
	}

	if err != nil {
		return nil, myerr.New(err.Error())
	}

	return tagI, nil
}
