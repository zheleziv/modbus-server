package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"zheleznovux.com/modbus-console/internal/configuration"
	"zheleznovux.com/modbus-console/internal/server"
	"zheleznovux.com/modbus-console/internal/win"
	myerr "zheleznovux.com/modbus-console/pkg"
)

func InitConfig(file string, s *server.Server) {
	conf, err := configuration.NewConfig(file)
	if err != nil {
		fmt.Printf("read config file err: %v\n", myerr.New(err.Error()))
		return
	}

	conf.AddObserver(s)
	s.Setup(conf)
}

func main() {
	server := server.New()

	if len(os.Args) > 1 {
		cmd := strings.ToLower(os.Args[1])
		if cmd == "sync" {
			server.Sync = true
		}
	}

	InitConfig("config.json", server)
	win.InitConfig("win_config.json")

	var wg sync.WaitGroup
	wg.Add(2)
	go server.Run()
	go win.Run(server)
	wg.Wait()
}
