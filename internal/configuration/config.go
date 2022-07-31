package configuration

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	myerr "zheleznovux.com/modbus-console/pkg"
)

// абстрактный класс, наследуемый ConfigurationDataApp и ConfigurationDataWin
type ConfigurationData interface {
	Setup(*ConfigHandler) error
}

type ConfigHandler struct {
	data           ConfigurationData // верифицированные данные из файла filename
	fileName       string            // имя файла с конифгурацией
	lastModifyTime int64             // время последней модификации файла, необходимо для горячей подгрузки конфигурации (reload())
	rwLock         sync.RWMutex      // синхронизация
	notifyList     []Notifyer        // массив классов обработчиков data
}

// обертка для выбора типа ConfigurationData и вызова Setup
func (c *ConfigHandler) parse() (ConfigurationData, error) {
	var tmpСonf ConfigurationData
	if strings.Contains(c.fileName, "win_") {
		tmpСonf = &ConfigurationDataWin{}
		if err := tmpСonf.(*ConfigurationDataWin).Setup(c); err != nil {
			return nil, myerr.New(err.Error())
		}
	} else {
		tmpСonf = &ConfigurationDataApp{}
		if err := tmpСonf.(*ConfigurationDataApp).Setup(c); err != nil {
			return nil, myerr.New(err.Error())
		}
	}

	return tmpСonf, nil
}

// конструктор
func NewConfig(fileName string) (conf *ConfigHandler, err error) {

	conf = &ConfigHandler{
		fileName: fileName,
	}

	m, err := conf.parse()
	if err != nil {
		fmt.Printf("parse conf error:%v\n", err)
		return nil, myerr.New(err.Error())
	}

	conf.rwLock.Lock()
	conf.data = m
	conf.rwLock.Unlock()

	go conf.reload()
	return conf, nil
}

func (c *ConfigHandler) GetConfig() ConfigurationData {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.data
}

func (c *ConfigHandler) reload() {
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		ticker.Stop()
		func() {
			f, err := os.Open(c.fileName)
			if err != nil {
				fmt.Printf("reload: open file error:%s\n", myerr.New(err.Error()))
				return
			}
			defer f.Close()

			fileInfo, err := f.Stat()
			if err != nil {
				fmt.Printf("stat file error:%s\n", myerr.New(err.Error()))
				return
			}

			curModifyTime := fileInfo.ModTime().Unix()
			if curModifyTime > c.lastModifyTime {
				m, err := c.parse()
				if err != nil {
					fmt.Printf("parse config error:%v\n", myerr.New(err.Error()))
					return
				}

				c.rwLock.Lock()
				c.data = m
				c.rwLock.Unlock()

				c.lastModifyTime = curModifyTime

				for _, n := range c.notifyList {
					n.Callback(c)
				}
			}

		}()
		ticker.Reset(time.Second * 5)
	}
}

// добавить смотрителя, реализующего класс Notifyer
func (c *ConfigHandler) AddObserver(n Notifyer) {
	c.notifyList = append(c.notifyList, n)
}
