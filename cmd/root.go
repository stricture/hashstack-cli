package cmd

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var authedSubCommands = []string{
	"project",
}

func topUsage(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

// RootCmd is the root level command for the cli.
var RootCmd = &cobra.Command{
	Use:   "hashstack-cli",
	Short: "Hashstack-cli is a cli client for Hashstack",
	Long:  "Hashstack-cli is a cli client for a remote Hashstack server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("try --help or hashstack-cli auth")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug(fmt.Sprintf("running sub command %s", cmd.Name()))
	},
}

var (
	flCfgFile   string
	flDebug     bool
	flServerURL string
	flToken     string
)

type config struct {
	ServerURL string `toml:"server_url"`
	Token     string `toml:"token"`
}

func initcfg() {
	if flCfgFile == "" {
		usr, err := user.Current()
		if err != nil || usr.HomeDir == "" {
			writeStdErrAndExit("There was an error locating your home directory.\nPlease use --config")
		}
		flCfgFile = filepath.Join(usr.HomeDir, "hashstack-cli.toml")
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

func init() {
	cobra.OnInitialize(initcfg)
	RootCmd.PersistentFlags().StringVar(&flCfgFile, "config", "", "config file (default is $HOME/hashstack-cli.toml)")
	RootCmd.PersistentFlags().BoolVar(&flDebug, "debug", false, "enable debug")
}

func debug(msg string) {
	if flDebug {
		fmt.Printf("DEBUG %s\n", msg)
	}
}

func ensureAuth(cmd *cobra.Command, args []string) {
	if flServerURL == "" || flToken == "" {
		writeStdErrAndExit("use hashstack-cli auth before continuing")
	}
}
