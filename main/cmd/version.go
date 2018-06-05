package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display kubespray-go version.",
	Long:  `Display kubespray-go version number.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.0")
	},
}