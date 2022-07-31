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

func (dwt *DWordTag) Setup(name string, address uint16, scanPeriod float64) error {
	var err error
	err = dwt.SetName(name)
	if err != nil {
		return myerr.New(err.Error())
	}
	err = dwt.SetAddress(address)
	if err != nil {
		return myerr.New(err.Error())
	}
	dwt.SetDataType()
	err = dwt.SetScanPeriod(scanPeriod)
	if err != nil {
		return myerr.New(err.Error())
	}
	dwt.SetState(false)
	return nil
}

func (dwt *DWordTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string
		DataType  string
		Address   uint16
		Value     uint32
		Timestamp string
		State     bool
	}{
		Name:      dwt.name,
		DataType:  dwt.dataType,
		Address:   dwt.address,
		Value:     dwt.value,
		Timestamp: dwt.timestamp,
		State:     dwt.state,
	})
}

//===================================Name
func (dwt *DWordTag) SetName(name string) error {
	tmp := strings.TrimSpace(name)
	if tmp == "" {
		return myerr.New("empty tag name")
	}
	dwt.name = tmp
	return nil
}
func (t *DWordTag) Name() string {
	return t.name
}

//===================================DataType
func (dwt *DWordTag) SetDataType() {
	dwt.dataType = DWORD_TYPE
}
func (dwt *DWordTag) DataType() string {
	return dwt.dataType
}

//===================================Address
func (dwt *DWordTag) Address() uint16 {
	return dwt.address
}
func (dwt *DWordTag) SetAddress(address uint16) error {
	if address >= UINT16_MAX_VALUE {
		return myerr.New("invalid tag address")
	}
	dwt.address = address
	return nil
}

//===================================TimeStamp
func (dwt *DWordTag) SetTimestamp() {
	now := time.Now()
	dwt.timestamp = now.Format(time.RFC3339)
}
func (dwt *DWordTag) Timestamp() string {
	return dwt.timestamp
}

//===================================State
func (dwt *DWordTag) SetState(state bool) {
	dwt.state = state
}
func (dwt *DWordTag) State() bool {
	return dwt.state
}

//===================================Value не интерфейсный метод
func (dwt *DWordTag) SetValue(value uint32) {
	dwt.rw.Lock()
	defer dwt.rw.Unlock()
	dwt.SetTimestamp()
	dwt.SetState(true)
	dwt.value = value
}
func (dwt *DWordTag) Value() uint32 {
	return dwt.value
}

//===================================ScanPeriod
func (dwt *DWordTag) ScanPeriod() float64 {
	return dwt.scanPeriod
}
func (dwt *DWordTag) SetScanPeriod(time float64) error {
	if time < 0 {
		return myerr.New("set scan period period < 0")
	}
	dwt.scanPeriod = time
	return nil
}
