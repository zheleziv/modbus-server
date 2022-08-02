package client

import (
	"sync"

	tag "zheleznovux.com/modbus-console/pkg/client/tag"
)

type ClientInterface interface {
	Start(stop chan struct{}, wg *sync.WaitGroup)

	Name() string
	Tags() []tag.TagInterface
	TagById(id int) (tag.TagInterface, error)
	TagByName(name string) (tag.TagInterface, error)
	SetTag(name string, address uint32, scanPeriod float64, dataType string) error
	// SetTags(tags []tag.TagInterface) error

	MarshalJSON() ([]byte, error) ///-----------------
}
