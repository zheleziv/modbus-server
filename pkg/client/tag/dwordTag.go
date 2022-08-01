package tag

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	myerr "zheleznovux.com/modbus-console/pkg"
)

type DWordTag struct {
	name       string
	dataType   string
	address    uint16
	scanPeriod float64
	value      uint32
	timestamp  string
	state      bool
	rw         sync.RWMutex
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
		DataType:  thisDWordTag.dataType,
		Address:   thisDWordTag.address,
		Value:     thisDWordTag.value,
		Timestamp: thisDWordTag.timestamp,
		State:     thisDWordTag.state,
	})
}

func (thisDWordTag *DWordTag) Setup(name string, address uint16, scanPeriod float64) error {
	var err error
	err = thisDWordTag.SetName(name)
	if err != nil {
		return myerr.New(err.Error())
	}
	err = thisDWordTag.SetAddress(address)
	if err != nil {
		return myerr.New(err.Error())
	}
	thisDWordTag.SetDataType()
	err = thisDWordTag.SetScanPeriod(scanPeriod)
	if err != nil {
		return myerr.New(err.Error())
	}
	thisDWordTag.SetState(false)
	return nil
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
func (thisDWordTag *DWordTag) Name() string {
	return thisDWordTag.name
}

//===================================DataType
func (thisDWordTag *DWordTag) SetDataType() {
	thisDWordTag.dataType = DWORD_TYPE
}
func (thisDWordTag *DWordTag) DataType() string {
	return thisDWordTag.dataType
}

//===================================Address
func (thisDWordTag *DWordTag) Address() uint16 {
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
func (thisDWordTag *DWordTag) Timestamp() string {
	return thisDWordTag.timestamp
}

//===================================State
func (thisDWordTag *DWordTag) SetState(state bool) {
	thisDWordTag.state = state
}
func (thisDWordTag *DWordTag) State() bool {
	return thisDWordTag.state
}

//===================================Value не интерфейсный метод
func (thisDWordTag *DWordTag) SetValue(value uint32) {
	thisDWordTag.rw.Lock()
	defer thisDWordTag.rw.Unlock()
	thisDWordTag.SetTimestamp()
	thisDWordTag.SetState(true)
	thisDWordTag.value = value
}
func (thisDWordTag *DWordTag) Value() uint32 {
	return thisDWordTag.value
}

//===================================ScanPeriod
func (thisDWordTag *DWordTag) ScanPeriod() float64 {
	return thisDWordTag.scanPeriod
}
func (thisDWordTag *DWordTag) SetScanPeriod(time float64) error {
	if time < 0 {
		return myerr.New("set scan period period < 0")
	}
	thisDWordTag.scanPeriod = time
	return nil
}
