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

func percentOf(current int, all int) float64 {
	percent := (float64(current) * float64(100)) / float64(all)
	return percent
}

func prettyUptime(uptime int64) string {
	var d, h, m, s int64
	var ret string
	if uptime <= 0 {
		uptime = 0
	}

	if uptime >= 86400 {
		d = uptime / 86400
		uptime = uptime - (86400 * d)
	}
	if uptime >= 3600 {
		h = uptime / 3600
		uptime = uptime - (3600 * h)
	}
	if uptime >= 60 {
		m = uptime / 60
		uptime = uptime - (60 * m)
	}
	s = uptime

	ret = fmt.Sprintf("%d day", d)
	if d != 1 {
		ret = fmt.Sprintf("%ss", ret)
	}
	ret = fmt.Sprintf("%s %d hour", ret, h)
	if h != 1 {
		ret = fmt.Sprintf("%ss", ret)
	}
	ret = fmt.Sprintf("%s %d minute", ret, m)
	if m != 1 {
		ret = fmt.Sprintf("%ss", ret)
	}
	ret = fmt.Sprintf("%s %d second", ret, s)
	if s != 1 {
		ret = fmt.Sprintf("%ss", ret)
	}
	return ret
}

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
		writeStdErrAndExit("Use hashstack-cli login before continuing.")
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
		debug("CONFIG: Could not decode configuration file")
		return
	}
	flServerURL = cfg.ServerURL
	flToken = cfg.Token
	debug(fmt.Sprintf("CONFIG: SERVER_URL - %s", flServerURL))
	debug(fmt.Sprintf("CONFIG: TOKEN - %s", flToken))
}

// writecfg will save the required data to the user's configuration file.
func writecfg() {
	tomlString := fmt.Sprintf("server_url = \"%s\"\ntoken = \"%s\"", flServerURL, flToken)
	os.Remove(flCfgFile)
	os.Mkdir(filepath.Dir(flCfgFile), 0655)
	fh, err := os.OpenFile(flCfgFile, os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit("There was an error opening the configuration file.")
	}
	defer fh.Close()
	if _, err := fh.WriteString(tomlString); err != nil {
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit("There was an error writing to the configuration file.")
	}
}

func initenv() {
	if flInsecure {
		debug("SECURITY: All requests are set to insecure")
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
		debug(fmt.Sprintf("Command: running sub command %s", cmd.Name()))
	},
}

func init() {
	cobra.OnInitialize(initcfg, initenv)
	RootCmd.PersistentFlags().StringVar(&flCfgFile, "config", "", "config file (default: $HOME/.hashstack/config)")
	RootCmd.PersistentFlags().BoolVar(&flInsecure, "insecure", false, "skip TLS certificate validation")
	RootCmd.PersistentFlags().BoolVar(&flDebug, "debug", false, "enable debug output")
}
