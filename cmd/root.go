package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "hashstack-cli",
	Short: "Hashstack-cli is a cli client for Hashstack",
	Long:  "Hashstack-cli is a cli client for a remote Hashstack server",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	//cobra.OnInitialize(initConfig)
}
