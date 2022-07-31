package tag

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	myerr "zheleznovux.com/modbus-console/pkg"
)

type WordTag struct {
	name       string
	dataType   string
	address    uint16
	scanPeriod float64
	value      uint16
	timestamp  string
	state      bool
	rw         sync.RWMutex
}

func (wt *WordTag) Setup(name string, address uint16, scanPeriod float64) error {
	var err error
	err = wt.SetName(name)
	if err != nil {
		return myerr.New(err.Error())
	}
	err = wt.SetAddress(address)
	if err != nil {
		return myerr.New(err.Error())
	}
	wt.SetDataType()
	err = wt.SetScanPeriod(scanPeriod)
	if err != nil {
		return myerr.New(err.Error())
	}
	wt.SetState(false)
	return nil
}

func (wt *WordTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string
		DataType  string
		Address   uint16
		Value     uint16
		Timestamp string
		State     bool
	}{
		Name:      wt.name,
		DataType:  wt.dataType,
		Address:   wt.address,
		Value:     wt.value,
		Timestamp: wt.timestamp,
		State:     wt.state,
	})
}

//===================================Name
func (t *WordTag) SetName(name string) error {
	tmp := strings.TrimSpace(name)
	if tmp == "" {
		return myerr.New("empty tag name")
	}
	t.name = tmp
	return nil
}
func (t *WordTag) Name() string {
	return t.name
}

//===================================DataType
func (t *WordTag) SetDataType() {
	t.dataType = WORD_TYPE
}
func (t *WordTag) DataType() string {
	return t.dataType
}

//===================================Address
func (t *WordTag) Address() uint16 {
	return t.address
}
func (t *WordTag) SetAddress(address uint16) error {
	if address >= UINT16_MAX_VALUE {
		return myerr.New("invalid tag address")

	}
	t.address = address
	return nil
}

//===================================TimeStamp
func (t *WordTag) SetTimestamp() {
	now := time.Now()
	t.timestamp = now.Format(time.RFC3339)
}
func (t *WordTag) Timestamp() string {
	return t.timestamp
}

//===================================State
func (t *WordTag) SetState(state bool) {
	t.state = state
}
func (t *WordTag) State() bool {
	return t.state
}

//===================================Value не интерфейсный метод
func (t *WordTag) SetValue(value uint16) {
	t.rw.Lock()
	defer t.rw.Unlock()
	t.SetTimestamp()
	t.SetState(true)
	t.value = value
}
func (t *WordTag) Value() uint16 {
	return t.value
}

//===================================ScanPeriod
func (t *WordTag) ScanPeriod() float64 {
	return t.scanPeriod
}
func (t *WordTag) SetScanPeriod(time float64) error {
	if time < 0 {
		return myerr.New("set scan period period < 0")
	}
	t.scanPeriod = time
	return nil
}
