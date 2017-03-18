package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Executes admin subcommands (-h or --help for more info)",
	Long: `
Executes a subcommand as an administrator (-h or --help for a list of subcommands.
    `,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

type impersonateRequest struct {
	Username string `json:"username"`
}

func getImpersonationToken(username string) string {
	body, err := postJSON("/api/admin/impersonate", impersonateRequest{
		Username: username,
	})
	if err != nil {
		writeStdErrAndExit(err.Error())
	}
	var response tokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit(new(invalidResponseError).Error())
	}
	return response.Token
}

var adminImpersonateCmd = &cobra.Command{
	Use:   "impersonate <username>",
	Short: "Create an authentication token for another user",
	Long: `
Create an authentication token for another user by username. This
can be useful when troubleshooting as another user.
    `,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("username is required")
			return
		}
		username := args[0]
		token := getImpersonationToken(username)
		fmt.Printf("Authentication token for %s: %s\n", username, token)
	},
}

func init() {
	adminCmd.AddCommand(adminImpersonateCmd)
	RootCmd.AddCommand(adminCmd)
}
