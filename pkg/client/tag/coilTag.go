package tag

import (
	"encoding/json"
	"strings"
	"time"

	myerr "zheleznovux.com/modbus-console/pkg"
)

type CoilTag struct {
	name       string
	address    uint16
	scanPeriod float64
	timestamp  string
	state      bool
	value      byte
}

// конструктор по умолчанию ???надо ли заполнять значениями по умолчанию???
func NewCoilTag() CoilTag {
	return CoilTag{}
}

// конструктор с параметрами
func NewCoilTagWithData(name string, address uint16, scanPeriod float64) (CoilTag, error) {
	thisCoilTag := NewCoilTag()
	err := thisCoilTag.SetName(name)
	if err != nil {
		return thisCoilTag, myerr.New(err.Error())
	}

	err = thisCoilTag.SetAddress(address)
	if err != nil {
		return thisCoilTag, myerr.New(err.Error())
	}

	err = thisCoilTag.SetScanPeriod(scanPeriod)
	if err != nil {
		return thisCoilTag, myerr.New(err.Error())
	}

	thisCoilTag.SetState(false)
	return thisCoilTag, nil
}

func (thisCoilTag *CoilTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string
		DataType  string
		Address   uint16
		Value     byte
		Timestamp string
		State     bool
	}{
		Name:      thisCoilTag.name,
		DataType:  thisCoilTag.DataType(),
		Address:   thisCoilTag.address,
		Value:     thisCoilTag.value,
		Timestamp: thisCoilTag.timestamp,
		State:     thisCoilTag.state,
	})
}

//===================================Name
func (thisCoilTag *CoilTag) SetName(name string) error {
	tmp := strings.TrimSpace(name)
	if tmp == "" {
		return myerr.New("empty tag name")
	}
	thisCoilTag.name = tmp
	return nil
}
func (thisCoilTag CoilTag) Name() string { // а если я не хочу давать укзатель на приватное поле класса? или как законстантить метод?
	return thisCoilTag.name
}

// //===================================DataType
func (thisCoilTag CoilTag) DataType() string {
	return COIL_TYPE
}

//===================================Address
func (thisCoilTag *CoilTag) SetAddress(address uint16) error {
	if address >= UINT16_MAX_VALUE {
		return myerr.New("invalid tag address")
	}
	thisCoilTag.address = address
	return nil
}
func (thisCoilTag CoilTag) Address() uint16 {
	return thisCoilTag.address
}

//===================================ScanPeriod
func (thisCoilTag *CoilTag) SetScanPeriod(time float64) error {
	if time < 0 {
		return myerr.New("set scan period < 0")
	}
	thisCoilTag.scanPeriod = time
	return nil
}
func (thisCoilTag CoilTag) ScanPeriod() float64 {
	return thisCoilTag.scanPeriod
}

//===================================Value
func (thisCoilTag *CoilTag) SetValue(value byte) {
	thisCoilTag.SetTimestamp()
	thisCoilTag.SetState(true)
	thisCoilTag.value = value
}
func (thisCoilTag CoilTag) Value() byte {
	return thisCoilTag.value
}

//===================================TimeStamp
func (thisCoilTag *CoilTag) SetTimestamp() {
	now := time.Now()
	thisCoilTag.timestamp = now.Format(time.RFC3339)
}
func (thisCoilTag CoilTag) Timestamp() string {
	return thisCoilTag.timestamp
}

//===================================DataState
func (thisCoilTag *CoilTag) SetState(state bool) {
	thisCoilTag.state = state
}
func (thisCoilTag CoilTag) State() bool {
	return thisCoilTag.state
}
