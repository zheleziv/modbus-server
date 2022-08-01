package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"zheleznovux.com/modbus-console/internal/configuration"

	myerr "zheleznovux.com/modbus-console/pkg"
	"zheleznovux.com/modbus-console/pkg/client"
	"zheleznovux.com/modbus-console/pkg/client/tag"
)

type Server struct {
	data []client.ClientInterface
	Sync bool
}

var changeCh chan int = make(chan int)

func (thisServer *Server) Callback(conf *configuration.ConfigHandler) {
	thisServer.Setup(conf)
	changeCh <- 1
}

func New() *Server {
	return &Server{}
}

func (thisServer *Server) Setup(confHandler *configuration.ConfigHandler) {
	config := confHandler.GetConfig()
	rtn := make([]client.ClientInterface, 0)

	tmpTN := config.(*configuration.ConfigurationDataApp)
	// проверка полученных данных
	// цикл по узлам
	for i := 0; i < len(tmpTN.NODES); i++ {
		k := 0
		j := i + 1
		// считаем количество неодинаковых имен
		for ; j < len(tmpTN.NODES); j++ {
			if strings.TrimSpace(tmpTN.NODES[i].Name) != strings.TrimSpace(tmpTN.NODES[j].Name) {
				k++
			}
		}
		// если все имена неодинаковые, то проверяем полученные данные
		// и добавляем в выходной массив новый узел
		if (j - i - 1) == k {
			var tmp client.ClientInterface
			nodes := tmpTN.NODES
			switch nodes[i].ConnectionType {
			case client.MODBUS_TCP:
				{
					var err error

					tmp, err = client.NewClinetModbus(nodes[i].IP, nodes[i].Port, nodes[i].ID, nodes[i].Name, nodes[i].Log, int(nodes[i].ConnectionAttempts), nodes[i].ConnectionTimeout)
					if err != nil {
						fmt.Println(myerr.New(err.Error()))
						continue
					}
					for j := range nodes[i].TAGS {
						err = tmp.SetTag(
							nodes[i].TAGS[j].Name,
							nodes[i].TAGS[j].Address,
							nodes[i].TAGS[j].ScanPeriod,
							nodes[i].TAGS[j].DataType)
						if err != nil {
							fmt.Println(myerr.New(err.Error()))
							continue
						}
					}
				}
			default:
				{
					fmt.Println("неизвестный тип подключения")
					continue
				}
			}
			rtn = append(rtn, tmp)
		}
	}

	thisServer.data = rtn
}

func (thisServer *Server) GetData() []client.ClientInterface {
	return thisServer.data
}

func (thisServer *Server) GetTagByName(name string) (tag.TagInterface, error) {

	split := strings.Split(name, ".")
	if len(split) != 2 {
		return nil, myerr.New("invalid name")
	}

	for i := range thisServer.data {
		if thisServer.data[i].Name() == split[0] {
			for j := range thisServer.data[i].Tags() {
				if thisServer.data[i].Tags()[j].Name() == split[1] {
					return thisServer.data[i].Tags()[j], nil
				}
			}
		}
	}
	return nil, myerr.New("no such name")
}

func (thisServer *Server) Save() {
	rankingsJson, err := json.Marshal(thisServer.data)
	if err != nil {
		fmt.Println(myerr.New(err.Error()))
		return
	}
	err = ioutil.WriteFile("output.json", rankingsJson, 0644)
	if err != nil {
		fmt.Println(myerr.New(err.Error()))
		return
	}
}

func (thisServer *Server) Run() {
	quit := make(chan struct{})
	var wg sync.WaitGroup

	saveTicker := time.NewTicker(10 * time.Second)
	for {
		// сигнал смены конфига
		select {
		case <-changeCh:
			{
				close(quit)
				wg.Wait()
				quit = make(chan struct{})

				for clientId := range thisServer.data {
					wg.Add(1)
					go thisServer.data[clientId].Start(quit, &wg)
				}
			}
		case <-saveTicker.C:
			{
				thisServer.Save()
			}
		}
	}
}
