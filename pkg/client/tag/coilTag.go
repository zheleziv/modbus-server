package tag

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	myerr "zheleznovux.com/modbus-console/pkg"
)

type CoilTag struct {
	name       string
	dataType   string
	address    uint16
	scanPeriod float64
	value      byte
	timestamp  string
	state      bool
	rw         sync.RWMutex
}

func (ct *CoilTag) Setup(name string, address uint16, scanPeriod float64) error {
	var err error
	err = ct.SetName(name)
	if err != nil {
		return myerr.New(err.Error())
	}
	err = ct.SetAddress(address)
	if err != nil {
		return myerr.New(err.Error())
	}
	ct.SetDataType()
	err = ct.SetScanPeriod(scanPeriod)
	if err != nil {
		return myerr.New(err.Error())
	}
	ct.SetState(false)
	return nil
}

func (ct *CoilTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string
		DataType  string
		Address   uint16
		Value     byte
		Timestamp string
		State     bool
	}{
		Name:      ct.name,
		DataType:  ct.dataType,
		Address:   ct.address,
		Value:     ct.value,
		Timestamp: ct.timestamp,
		State:     ct.state,
	})
}

//===================================Name
func (t *CoilTag) SetName(name string) error {
	tmp := strings.TrimSpace(name)
	if tmp == "" {
		return myerr.New("empty tag name")
	}
	t.name = tmp
	return nil
}
func (t *CoilTag) Name() string {
	return t.name
}

//===================================DataType
func (t *CoilTag) SetDataType() {
	t.dataType = COIL_TYPE
}
func (t *CoilTag) DataType() string {
	return t.dataType
}

//===================================Address
func (t *CoilTag) SetAddress(address uint16) error {
	if address >= UINT16_MAX_VALUE {
		return myerr.New("invalid tag address")
	}
	t.address = address
	return nil
}
func (t *CoilTag) Address() uint16 {
	return t.address
}

//===================================ScanPeriod
func (t *CoilTag) SetScanPeriod(time float64) error {
	if time < 0 {
		return myerr.New("set scan period < 0")
	}
	t.scanPeriod = time
	return nil
}
func (t *CoilTag) ScanPeriod() float64 {
	t.rw.Lock()
	defer t.rw.Unlock()
	return t.scanPeriod
}

//===================================Value
func (t *CoilTag) SetValue(value byte) {
	t.rw.Lock()
	defer t.rw.Unlock()
	t.SetTimestamp()
	t.SetState(true)
	t.value = value
}
func (t *CoilTag) Value() byte {
	return t.value
}

//===================================TimeStamp
func (t *CoilTag) SetTimestamp() {
	now := time.Now()
	t.timestamp = now.Format(time.RFC3339)
}
func (t *CoilTag) Timestamp() string {
	return t.timestamp
}

//===================================DataState
func (t *CoilTag) SetState(state bool) {
	t.state = state
}
func (t *CoilTag) State() bool {
	return t.state
}
