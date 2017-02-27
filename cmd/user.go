package cmd

import (
	"fmt"
	"time"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
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
		var users []hashstack.User
		if err := getRangeJSON("/api/users", &users); err != nil {
			writeStdErrAndExit(err.Error())
		}
		tbl := uitable.New()
		tbl.AddRow("ID", "USERNAME", "UPDATED")
		for _, u := range users {
			tm := time.Unix(u.UpdatedAt, 0)
			tbl.AddRow(u.ID, u.Username, tm)
		}
		fmt.Println(tbl)
	},
}

func init() {
	RootCmd.AddCommand(userCmd)
}
