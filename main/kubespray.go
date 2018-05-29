package main

import (
	"github.com/hekonsek/kubespray-go"
	"fmt"
)

func main() {
	kubespray := kubespray_go.NewKubespray()
	err := kubespray.Provision()
	if err != nil {
		fmt.Println(err)
	}
}