package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and exit",
	Long:  "Print the version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s", version)
	},
}
