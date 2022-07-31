package configuration

import (
	"encoding/json"
	"os"

	myerr "zheleznovux.com/modbus-console/pkg"
)

type ConfigurationDataApp struct {
	NODES []Node
}

type Node struct {
	Name               string
	ConnectionType     string
	IP                 string
	Port               int
	ID                 uint8
	Log                bool
	ConnectionTimeout  float64
	ConnectionAttempts uint
	TAGS               []Tag
}

type Tag struct {
	Name       string
	Address    uint32
	DataType   string
	ScanPeriod float64
}

func (tn *ConfigurationDataApp) Setup(c *ConfigHandler) error {
	content, err := os.ReadFile(c.fileName)
	if err != nil {
		return myerr.New(err.Error())
	}
	var tmpTN ConfigurationDataApp
	err = json.Unmarshal(content, &tmpTN)
	if err != nil {
		return myerr.New(err.Error())
	}
	*tn = tmpTN
	return nil
}
