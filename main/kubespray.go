package main

import (
	"github.com/hekonsek/kubespray-go"
	"fmt"
	"os"
)

func main() {
	kubespray, err := kubespray_go.NewKubespray(os.Args[1])
	kubespray.AnsibleBecome = true
	kubespray.AnsibleUser = "fedora"
	if err != nil {
		fmt.Println(err)
	}
	err = kubespray.Provision()
	if err != nil {
		fmt.Println(err)
	}
}