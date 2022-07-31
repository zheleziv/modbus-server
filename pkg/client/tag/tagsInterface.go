package tag

import (
	myerr "zheleznovux.com/modbus-console/pkg"
)

type TagInterface interface {
	SetName(string) error
	Name() string

	SetAddress(uint16) error
	Address() uint16

	SetScanPeriod(float64) error
	ScanPeriod() float64

	SetDataType()
	DataType() string

	SetState(bool)
	State() bool

	SetTimestamp()
	Timestamp() string

	Setup(name string, address uint16, scanPeriod float64) error
}

func NewTag(name string, address uint16, scanPeriod float64, dataType string) (TagInterface, error) {
	var tagI TagInterface
	if dataType == "" {
		return nil, myerr.New("invalid dataType")
	}
	switch dataType {
	case COIL_TYPE:
		{
			var tag CoilTag
			tagI = &tag
		}
	case WORD_TYPE:
		{
			var tag WordTag
			tagI = &tag
		}
	case DWORD_TYPE:
		{
			var tag DWordTag
			tagI = &tag
		}
	default:
		return nil, myerr.New("invalid dataType")
	}

	err := tagI.Setup(name, address, scanPeriod)
	if err != nil {
		return nil, myerr.New(err.Error())
	}
	return tagI, nil
}
