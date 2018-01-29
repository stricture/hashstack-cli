package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var version = "2.1.0"

type serverVersion struct {
	Version string `json:"version"`
}

func getServerVersion() (string, error) {
	debug(fmt.Sprintf("HTTP: GET %s", "/version"))

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", flServerURL, "/version"), nil)
	if err != nil {
		debug(fmt.Sprintf("Error: %s", err.Error()))
		return "", new(requestCreateError)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug(fmt.Sprintf("Error: %s", err.Error()))
		if strings.Contains(err.Error(), "x509") {
			return "", new(invalidCertError)
		}
		return "", new(requestError)
	}
	if err := respToError(resp); err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		debug(fmt.Sprintf("Error: %s", err.Error()))
		return "", new(invalidResponseError)
	}
	var version serverVersion
	if err := json.Unmarshal(body, &version); err != nil {
		debug(fmt.Sprintf("Error: %s", err.Error()))
		return "", new(jsonServerError)
	}
	return version.Version, nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print client and server version and exit.",
	Long:  "Print client and server version and exit.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Client Version: %s\n", version)
		serverv, err := getServerVersion()
		if err != nil {
			writeStdErrAndExit(err.Error())
			return
		}
		fmt.Printf("Server Version: %s\n", serverv)
	},
}
