package cmd

import (
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

type updateSelf struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

var passwdCmd = &cobra.Command{
	Use:   "passwd",
	Short: "Change your current password.",
	Long: `
Change your current password.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		var user hashstack.User
		if err := getJSON("/api/users/self", &user); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Printf("Current password: ")
		currentpass, err := gopass.GetPasswdMasked()
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error reading your password.")
		}

		fmt.Printf("New password: ")
		pass, err := gopass.GetPasswdMasked()
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error reading your password.")
		}
		fmt.Printf("Confirm password: ")
		confirmpass, err := gopass.GetPasswdMasked()
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error reading your password.")
		}

		if string(pass) != string(confirmpass) {
			writeStdErrAndExit("Passwords did not match")
		}

		if _, err := postJSON("/api/users/self", &updateSelf{
			Username:    user.Username,
			Password:    string(currentpass),
			NewPassword: string(pass),
		}); err != nil {
			writeStdErrAndExit(err.Error())
		}

		fmt.Println("\nPassword updated. Please login again.")
	},
}

func init() {
	RootCmd.AddCommand(passwdCmd)
}
