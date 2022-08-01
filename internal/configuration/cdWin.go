package configuration

import (
	"encoding/json"
	"os"

	myerr "zheleznovux.com/modbus-console/pkg"
)

type ConfigurationDataWin struct {
	NODES []NodeTag
}

type NodeTag struct {
	Name           string
	Log            bool
	StateCondition string
	ValueCondition string
	Logic          string
	Action         string
	ActionTimeout  float64
	ScanPeriod     float64
}

func (thisConfig *ConfigurationDataWin) Setup(c *ConfigHandler) error {
	content, err := os.ReadFile(c.fileName)
	if err != nil {
		return myerr.New(err.Error())
	}

	var tmpTN ConfigurationDataWin
	err = json.Unmarshal(content, &tmpTN)
	if err != nil {
		return myerr.New(err.Error())
	}

	// проверка на одинаковые имена
	for i := 0; i < len(tmpTN.NODES); i++ {
		k := 0
		j := i + 1
		for ; j < len(tmpTN.NODES); j++ {
			if tmpTN.NODES[i].Name != tmpTN.NODES[j].Name {
				k++
			}
		}

		if (j - i - 1) == k {
			// если прошел проверку добавляем
			thisConfig.NODES = append(thisConfig.NODES, tmpTN.NODES[i])
		}
	}

	return nil
}
