package client

import (
	"sync"

	myerr "zheleznovux.com/modbus-console/pkg"
	tag "zheleznovux.com/modbus-console/pkg/client/tag"
)

type ClientInterface interface {
	Start(stop chan struct{}, wg *sync.WaitGroup)

	Name() string
	Type() string

	Tags() []tag.TagInterface
	TagById(id int) (tag.TagInterface, error)
	TagByName(name string) (tag.TagInterface, error)

	SetTag(name string, address uint32, scanPeriod float64, dataType string) error

	MarshalJSON() ([]byte, error)
}

func New(connectionType string, ip string, port int, slaveID uint8, name string, debug bool, ConnectionAttempts uint, ConnectionTimeout float64) (ClientInterface, error) {
	switch connectionType {
	case MODBUS_TCP:
		tmp, err := NewClinetModbus(
			ip,
			port,
			slaveID,
			name,
			debug,
			ConnectionAttempts,
			ConnectionTimeout)

		if err != nil {
			return tmp, myerr.New(err.Error())
		}
		return tmp, nil
	default:
		return nil, myerr.New("неизвестный тип подключения клиента")
	}
}
