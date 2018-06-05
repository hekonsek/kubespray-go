package cmd

import (
	"github.com/spf13/cobra"
	"github.com/hekonsek/kubespray-go"
	"os"
	"fmt"
)

var verbose bool
var ansibleBecome bool
var ansibleUser string

func init() {
	provisionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enables verbose command execution mode.")
	provisionCmd.Flags().BoolVarP(&ansibleBecome, "ansible-become", "b", false, "Enables Ansible become mode.")
	provisionCmd.Flags().StringVarP(&ansibleUser, "ansible-user", "u", "root", "Specified Ansible user.")
	rootCmd.AddCommand(provisionCmd)
}

var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision Kubernetes cluster using kubespray.",
	Long:  `Provision Kubernetes cluster using kubespray.`,
	Run: func(cmd *cobra.Command, args []string) {
			kubespray, err := kubespray_go.NewKubespray(os.Args[1])
			kubespray.Verbose = verbose
			kubespray.AnsibleBecome = ansibleBecome
			kubespray.AnsibleUser = ansibleUser
			if err != nil {
				fmt.Println(err)
				return
			}
			err = kubespray.Provision()
			if err != nil {
				fmt.Println(err)
				return
			}
	},
}