package win

import (
	"fmt"
	"sync"
	"sync/atomic"

	configuration "zheleznovux.com/modbus-console/internal/configuration"
	server "zheleznovux.com/modbus-console/internal/server"
	commander "zheleznovux.com/modbus-console/internal/win/commander"
	myerr "zheleznovux.com/modbus-console/pkg"
)

type WinNotifyerApp struct {
}

type WinConfig struct {
	nodeCommand []configuration.NodeTag
}

type WinConfigMgr struct {
	config atomic.Value
}

var winConfigMgr = &WinConfigMgr{}
var changeCh chan int = make(chan int)

func (a *WinNotifyerApp) Callback(conf *configuration.ConfigHandler) {
	changeCh <- 1
}

func (a *WinConfigMgr) Callback(conf *configuration.ConfigHandler) {
	winConfig := &WinConfig{}
	winConfig.nodeCommand = conf.GetConfig().(*configuration.ConfigurationDataWin).NODES
	winConfigMgr.config.Store(winConfig)
	changeCh <- 1
}

func InitConfig(file string) {
	conf, err := configuration.NewConfig(file)
	if err != nil {
		fmt.Printf("read config file err: %v\n", myerr.New(err.Error()))
		return
	}

	conf.AddObserver(winConfigMgr)
}

func Run(th *server.Server) {
	quit := make(chan struct{})
	var wg sync.WaitGroup

	for {
		<-changeCh
		close(quit)
		wg.Wait()
		quit = make(chan struct{})

		winConfig := winConfigMgr.config.Load().(*WinConfig)
		for i := range winConfig.nodeCommand {
			var com commander.Commander

			tag, err := th.GetTagByName(winConfig.nodeCommand[i].Name)
			if err != nil {
				fmt.Println(myerr.New(err.Error()))
				continue
			}
			com.Setup(winConfig.nodeCommand[i], &tag)
			wg.Add(1)
			go com.StartChecking(quit, &wg)
		}

	}

}
