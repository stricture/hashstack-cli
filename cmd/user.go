package cmd

import (
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "users",
	Short: "Displays a list of hashstack users",
	Long: `
You may list users here. Interactions with hashstack for other operations
should be conducted directly on the server by an administrator.

Users can be added or removed from a project using the projects command.
    `,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(userCmd)
}
