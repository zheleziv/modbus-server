package tag

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	myerr "zheleznovux.com/modbus-console/pkg"
)

type DWordTag struct {
	name       string
	address    uint16
	scanPeriod float64
	value      uint32
	timestamp  string
	state      bool
	ReadFunc   func() (uint32, error)
}

func (thisDWordTag *DWordTag) ReadDevice() error {
	val, err := thisDWordTag.ReadFunc()
	if err != nil {
		thisDWordTag.SetState(false)
		return myerr.New(err.Error())
	}

	thisDWordTag.setValue(val)
	return nil
}

// конструктор по умолчанию
func NewDWordTag() *DWordTag {
	return &DWordTag{}
}

// конструктор с параметрами
func NewDWordTagWithData(name string, address uint16, scanPeriod float64) (*DWordTag, error) {
	thisDWordTag := NewDWordTag()

	err := thisDWordTag.SetName(name)
	if err != nil {
		return thisDWordTag, myerr.New(err.Error())
	}

	err = thisDWordTag.SetAddress(address)
	if err != nil {
		return thisDWordTag, myerr.New(err.Error())
	}

	err = thisDWordTag.SetScanPeriod(scanPeriod)
	if err != nil {
		return thisDWordTag, myerr.New(err.Error())
	}
	fmt.Println("213")
	thisDWordTag.SetState(false)
	return thisDWordTag, nil
}

func (thisDWordTag *DWordTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string
		DataType  string
		Address   uint16
		Value     uint32
		Timestamp string
		State     bool
	}{
		Name:      thisDWordTag.name,
		DataType:  thisDWordTag.DataType(),
		Address:   thisDWordTag.address,
		Value:     thisDWordTag.value,
		Timestamp: thisDWordTag.timestamp,
		State:     thisDWordTag.state,
	})
}

//===================================Name
func (thisDWordTag *DWordTag) SetName(name string) error {
	tmp := strings.TrimSpace(name)
	if tmp == "" {
		return myerr.New("empty tag name")
	}
	thisDWordTag.name = tmp
	return nil
}
func (thisDWordTag DWordTag) Name() string {
	return thisDWordTag.name
}

//===================================DataType
func (thisDWordTag DWordTag) DataType() string {
	return DWORD_TYPE
}

//===================================Address
func (thisDWordTag DWordTag) Address() uint16 {
	return thisDWordTag.address
}
func (thisDWordTag *DWordTag) SetAddress(address uint16) error {
	if address >= UINT16_MAX_VALUE {
		return myerr.New("invalid tag address")
	}
	thisDWordTag.address = address
	return nil
}

//===================================TimeStamp
func (thisDWordTag *DWordTag) SetTimestamp() {
	now := time.Now()
	thisDWordTag.timestamp = now.Format(time.RFC3339)
}
func (thisDWordTag DWordTag) Timestamp() string {
	return thisDWordTag.timestamp
}

//===================================State
func (thisDWordTag *DWordTag) SetState(state bool) {
	thisDWordTag.state = state
}
func (thisDWordTag DWordTag) State() bool {
	return thisDWordTag.state
}

//===================================Value не интерфейсный метод
func (thisDWordTag *DWordTag) setValue(value uint32) {
	thisDWordTag.SetTimestamp()
	thisDWordTag.SetState(true)
	thisDWordTag.value = value
}
func (thisDWordTag DWordTag) Value() uint32 {
	return thisDWordTag.value
}

//===================================ScanPeriod
func (thisDWordTag DWordTag) ScanPeriod() float64 {
	return thisDWordTag.scanPeriod
}
func (thisDWordTag *DWordTag) SetScanPeriod(time float64) error {
	if time < 0 {
		return myerr.New("set scan period period < 0")
	}
	thisDWordTag.scanPeriod = time
	return nil
}
