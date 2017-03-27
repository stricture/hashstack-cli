package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"net/url"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

type tokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

var loginCmd = &cobra.Command{
	Use:   "login <server_url> <username>",
	Short: "Login and cache a session token for the remote server.",
	Long: `
This command will prompt for your password and send it along with your username to the server at server_url.
The token returned along with the server_url will be saved in your home directory for all additional requests.
    `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			writeStdErrAndExit("server_url and username are required.")
		}
		serverURL := args[0]
		username := args[1]

		if _, err := url.Parse(serverURL); err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("The provided URL is not valid. It is possible that you provided the arguments our of order.")
		}

		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswdMasked()
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error reading your password.")
		}

		data, err := json.Marshal(tokenRequest{
			Username: username,
			Password: string(pass),
		})
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error preparing your request for the server")
		}
		path := fmt.Sprintf("%s/token", serverURL)
		debug(fmt.Sprintf("HTTP: POST %s", path))
		resp, err := http.Post(path, "application/json", bytes.NewBuffer(data))
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit(new(requestError).Error())
		}
		var response tokenResponse
		switch resp.StatusCode {
		case 401:
			writeStdErrAndExit(new(authError).Error())
		case 500:
			writeStdErrAndExit(new(internalServerError).Error())
		case 201:
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				debug(fmt.Sprintf("Error: %s", err.Error()))
				writeStdErrAndExit(new(invalidResponseError).Error())
			}
		default:
			debug(fmt.Sprintf("HTTP: unexpected response code - %d", resp.StatusCode))
			writeStdErrAndExit("There was an unexpected response from the server.")
		}
		flServerURL = serverURL
		flToken = response.Token
		writecfg()
		fmt.Printf("Authentication credentials cached in %s.\n", flCfgFile)
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
