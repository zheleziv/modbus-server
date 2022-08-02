package tag

import (
	myerr "zheleznovux.com/modbus-console/pkg"
)

type TagInterface interface { // виртуальные поля, которыми должны обладать все тэги
	SetName(string) error // имя
	Name() string

	SetAddress(uint16) error // адрес
	Address() uint16

	SetScanPeriod(float64) error // период сканирования
	ScanPeriod() float64

	// SetDataType()
	DataType() string // строка указывающая на тип данных тэга

	SetState(bool) // состоаяние подключения
	State() bool

	SetTimestamp() // временная ветка
	Timestamp() string
}

func NewTag(name string, address uint16, scanPeriod float64, dataType string) (TagInterface, error) {
	var tagI TagInterface
	var err error

	if dataType == "" {
		return nil, myerr.New("invalid dataType")
	}
	switch dataType {
	case COIL_TYPE:
		{
			var tag CoilTag
			tag, err = NewCoilTagWithData(name, address, scanPeriod)
			tagI = &tag
		}
	case WORD_TYPE:
		{
			var tag WordTag
			tag, err = NewWordTagWithData(name, address, scanPeriod)
			tagI = &tag
		}
	case DWORD_TYPE:
		{
			var tag DWordTag
			tag, err = NewDWordTagWithData(name, address, scanPeriod)
			tagI = &tag
		}
	default:
		return nil, myerr.New("invalid dataType")
	}

	if err != nil {
		return nil, myerr.New(err.Error())
	}

	return tagI, nil
}
