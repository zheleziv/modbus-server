package tag

import (
	"encoding/json"
	"strings"
	"time"

	myerr "zheleznovux.com/modbus-console/pkg"
)

type WordTag struct {
	name       string
	address    uint16
	scanPeriod float64
	value      uint16
	timestamp  string
	state      bool
}

// конструктор по умолчанию
func NewWordTag() WordTag {
	return WordTag{}
}

// конструктор с параметрами
func NewWordTagWithData(name string, address uint16, scanPeriod float64) (WordTag, error) {
	thisWordTag := NewWordTag()
	err := thisWordTag.SetName(name)
	if err != nil {
		return thisWordTag, myerr.New(err.Error())
	}

	err = thisWordTag.SetAddress(address)
	if err != nil {
		return thisWordTag, myerr.New(err.Error())
	}

	err = thisWordTag.SetScanPeriod(scanPeriod)
	if err != nil {
		return thisWordTag, myerr.New(err.Error())
	}

	thisWordTag.SetState(false)
	return thisWordTag, nil
}

func (thisWordTag *WordTag) Setup(name string, address uint16, scanPeriod float64) error {
	var err error
	err = thisWordTag.SetName(name)
	if err != nil {
		return myerr.New(err.Error())
	}
	err = thisWordTag.SetAddress(address)
	if err != nil {
		return myerr.New(err.Error())
	}
	err = thisWordTag.SetScanPeriod(scanPeriod)
	if err != nil {
		return myerr.New(err.Error())
	}
	thisWordTag.SetState(false)
	return nil
}

func (thisWordTag *WordTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string
		DataType  string
		Address   uint16
		Value     uint16
		Timestamp string
		State     bool
	}{
		Name:      thisWordTag.name,
		DataType:  thisWordTag.DataType(),
		Address:   thisWordTag.address,
		Value:     thisWordTag.value,
		Timestamp: thisWordTag.timestamp,
		State:     thisWordTag.state,
	})
}

//===================================Name
func (thisWordTag *WordTag) SetName(name string) error {
	tmp := strings.TrimSpace(name)
	if tmp == "" {
		return myerr.New("empty tag name")
	}
	thisWordTag.name = tmp
	return nil
}
func (thisWordTag *WordTag) Name() string {
	return thisWordTag.name
}

func (thisWordTag *WordTag) DataType() string {
	return WORD_TYPE
}

//===================================Address
func (thisWordTag *WordTag) Address() uint16 {
	return thisWordTag.address
}
func (thisWordTag *WordTag) SetAddress(address uint16) error {
	if address >= UINT16_MAX_VALUE {
		return myerr.New("invalid tag address")

	}
	thisWordTag.address = address
	return nil
}

//===================================TimeStamp
func (thisWordTag *WordTag) SetTimestamp() {
	now := time.Now()
	thisWordTag.timestamp = now.Format(time.RFC3339)
}
func (thisWordTag *WordTag) Timestamp() string {
	return thisWordTag.timestamp
}

//===================================State
func (thisWordTag *WordTag) SetState(state bool) {
	thisWordTag.state = state
}
func (thisWordTag *WordTag) State() bool {
	return thisWordTag.state
}

//===================================Value не интерфейсный метод
func (thisWordTag *WordTag) SetValue(value uint16) {
	// thisWordTag.rw.Lock()
	// defer thisWordTag.rw.Unlock()
	thisWordTag.SetTimestamp()
	thisWordTag.SetState(true)
	thisWordTag.value = value
}
func (thisWordTag *WordTag) Value() uint16 {
	return thisWordTag.value
}

//===================================ScanPeriod
func (thisWordTag *WordTag) ScanPeriod() float64 {
	return thisWordTag.scanPeriod
}
func (thisWordTag *WordTag) SetScanPeriod(time float64) error {
	if time < 0 {
		return myerr.New("set scan period period < 0")
	}
	thisWordTag.scanPeriod = time
	return nil
}
