package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

func displayUser(user hashstack.User) {
	fmt.Printf("ID.......: %d\n", user.ID)
	fmt.Printf("Username.: %s\n", user.Username)
	fmt.Printf("-----------------------\n")
}

func displayUsers(users []hashstack.User) {
	for _, u := range users {
		displayUser(u)
	}
}

func getUsers() []hashstack.User {
	var users []hashstack.User
	if err := getRangeJSON("/api/users", &users); err != nil {
		writeStdErrAndExit(err.Error())
	}
	return users
}

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
		users := getUsers()
		displayUsers(users)
	},
}

func init() {
	RootCmd.AddCommand(userCmd)
}
