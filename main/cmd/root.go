package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubespray-go",
	Short: "kubespray-go is a command-line wrapper around kubespray-go GoLang library.",
	Long: `kubespray-go is a command-line wrapper around kubespray-go GoLang library.
                See https://github.com/hekonsek/kubespray-go`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}