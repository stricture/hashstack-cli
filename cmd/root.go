package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

// Global command line flags.
var (
	flCfgFile   string
	flInsecure  bool
	flDebug     bool
	flServerURL string
	flToken     string
)

type config struct {
	ServerURL string `toml:"server_url"`
	Token     string `toml:"token"`
}

func debug(msg string) {
	if flDebug {
		fmt.Printf("DEBUG: %s\n", msg)
	}
}

var authedSubCommands = []string{
	"project",
}

func ensureAuth(cmd *cobra.Command, args []string) {
	if flServerURL == "" || flToken == "" {
		writeStdErrAndExit("Use hashstack-cli login before continuing")
	}
}

// initcfg will load the configurationfile in the user's home directory.
func initcfg() {
	if flCfgFile == "" {
		usr, err := user.Current()
		if err != nil || usr.HomeDir == "" {
			writeStdErrAndExit("There was an error locating your home directory.\nPlease use --config")
		}
		flCfgFile = filepath.Join(usr.HomeDir, ".hashstack", "config")
	}
	debug(fmt.Sprintf("configuration file: %s", flCfgFile))
	var cfg config
	if _, err := toml.DecodeFile(flCfgFile, &cfg); err != nil {
		debug("no configuration file")
		return
	}
	flServerURL = cfg.ServerURL
	flToken = cfg.Token
	debug(fmt.Sprintf("server url: %s", flServerURL))
	debug(fmt.Sprintf("token: %s", flToken))
}

// writecfg will save the required data to the user's configuration file.
func writecfg() {
	tomlString := fmt.Sprintf("server_url = \"%s\"\ntoken = \"%s\"", flServerURL, flToken)
	os.Remove(flCfgFile)
	os.Mkdir(filepath.Dir(flCfgFile), 0655)
	fh, err := os.OpenFile(flCfgFile, os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		debug(err.Error())
		writeStdErrAndExit("There was an error opening the configuration file.")
	}
	defer fh.Close()
	if _, err := fh.WriteString(tomlString); err != nil {
		debug(err.Error())
		writeStdErrAndExit("There was an error writing to the configuration file.")
	}
}

func initenv() {
	if flInsecure {
		debug("setting transport to insecure")
		http.DefaultTransport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
}

// RootCmd is the root level command for the cli.
// Executing this command will print the usage information and exit.
var RootCmd = &cobra.Command{
	Use:   "hashstack-cli",
	Short: "Execute commands against a Hashstack server. Try -h or --help for more information.",
	Long:  "Execute commands against a Hashstack server. Try -h or --help for more information.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug(fmt.Sprintf("running sub command %s", cmd.Name()))
	},
}

func init() {
	cobra.OnInitialize(initcfg, initenv)
	RootCmd.PersistentFlags().StringVar(&flCfgFile, "config", "", "config file (default: $HOME/.hashstack/config)")
	RootCmd.PersistentFlags().BoolVarP(&flInsecure, "insecure", "k", false, "skip TLS certificate validation")
	RootCmd.PersistentFlags().BoolVar(&flDebug, "debug", false, "enable debug output")
}
