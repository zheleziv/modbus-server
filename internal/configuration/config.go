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
func (thisConfigHandler *ConfigHandler) parse() (ConfigurationData, error) {
	var tmpСonf ConfigurationData
	if strings.Contains(thisConfigHandler.fileName, "win_") {
		tmpСonf = &ConfigurationDataWin{}
		if err := tmpСonf.(*ConfigurationDataWin).Setup(thisConfigHandler); err != nil {
			return nil, myerr.New(err.Error())
		}
	} else {
		tmpСonf = &ConfigurationDataApp{}
		if err := tmpСonf.(*ConfigurationDataApp).Setup(thisConfigHandler); err != nil {
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

func (thisConfigHandler *ConfigHandler) GetConfig() ConfigurationData {
	thisConfigHandler.rwLock.RLock()
	defer thisConfigHandler.rwLock.RUnlock()
	return thisConfigHandler.data
}

func (thisConfigHandler *ConfigHandler) reload() {
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		ticker.Stop()
		func() {
			f, err := os.Open(thisConfigHandler.fileName)
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
			if curModifyTime > thisConfigHandler.lastModifyTime {
				m, err := thisConfigHandler.parse()
				if err != nil {
					fmt.Printf("parse config error:%v\n", myerr.New(err.Error()))
					return
				}

				thisConfigHandler.rwLock.Lock()
				thisConfigHandler.data = m
				thisConfigHandler.rwLock.Unlock()

				thisConfigHandler.lastModifyTime = curModifyTime

				for _, n := range thisConfigHandler.notifyList {
					n.Callback(thisConfigHandler)
				}
			}
		}()
		ticker.Reset(time.Second * 5)
	}
}

// добавить смотрителя, реализующего класс Notifyer
func (thisConfigHandler *ConfigHandler) AddObserver(n Notifyer) {
	thisConfigHandler.notifyList = append(thisConfigHandler.notifyList, n)
}
