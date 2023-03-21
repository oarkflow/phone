package cmd

import (
	"fmt"
	"github.com/oarkflow/phone"
)

func Execute() {
	initBackgroundServices()
}

func initBackgroundServices() {
	if err := phone.LoadNetworks(); err != nil {
		fmt.Println(err)
	}
}
