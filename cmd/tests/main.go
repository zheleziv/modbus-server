package main

import (
	"fmt"

	myerr "zheleznovux.com/modbus-console/pkg"
	"zheleznovux.com/modbus-console/pkg/client"
)

func main() {
	cc, err := client.NewClinetModbus("127.0.0.1", 432, 1, "tagser", false, 2, 132)
	if err != nil {
		fmt.Println(myerr.New(err.Error()))
	}
	t, err := cc.TagById(13123123)
	if err != nil {
		fmt.Println(myerr.New(err.Error()))
	} else {
		fmt.Println(t)
	}

}
