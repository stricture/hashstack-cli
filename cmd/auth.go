package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(authCmd)
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

var authCmd = &cobra.Command{
	Use:   "auth [server_url] [username]",
	Short: "Authenticate to Hashstack server",
	Long: `
auth will prompt for your password and send it along with your username
to the server at server_url. The token returned along with the server_url
will be saved in your home directory for all additional requests.
    `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			writeStdErrAndExit("server_url and username are required")
		}
		serverURL := args[0]
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswdMasked()
		if err != nil {
			writeStdErrAndExit("invalid password provided")
		}
		data, err := json.Marshal(authRequest{
			Username: args[1],
			Password: string(pass),
		})
		if err != nil {
			writeStdErrAndExit("invalid username or password")
		}
		resp, err := http.Post(fmt.Sprintf("%s/token", serverURL), "application/json", bytes.NewBuffer(data))
		if err != nil {
			writeStdErrAndExit("error sending request")
		}
		var response authResponse
		switch resp.StatusCode {
		case 401:
			writeStdErrAndExit("invalid username or password")
		case 500:
			writeStdErrAndExit("internal server error")
		case 201:
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				writeStdErrAndExit("error decoding response")
			}
		default:
			writeStdErrAndExit("unexpected response from server")
		}
		tomlString := fmt.Sprintf("server_url = \"%s\"\ntoken = \"%s\"", serverURL, response.Token)
		usr, err := user.Current()
		if err != nil || usr.HomeDir == "" {
			fmt.Println("Your home directory could not be located!")
			fmt.Println("Please save the following into a configuration file:")
			fmt.Println(tomlString)
		}
		filename := filepath.Join(usr.HomeDir, "hashstack-cli.toml")
		fh, err := os.OpenFile(filename, os.O_CREATE, 0655)
		if err != nil {
			writeStdErrAndExit("error opening file")
		}
		defer fh.Close()
		fh.WriteString(tomlString)
		fmt.Println("Authentication information saved. You can now run additional commands")
	},
}
