package cmd

import (
	"fmt"

	"sort"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

var modeCmd = &cobra.Command{
	Use:    "modes",
	Short:  "Prints a list of supported hash modes",
	Long:   "Prints a list of supported hash modes",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		var hashModes []hashstack.HashMode
		if err := getJSON("/api/hash_modes", &hashModes); err != nil {
			writeStdErrAndExit(err.Error())
		}
		sort.Slice(hashModes, func(i, j int) bool {
			return hashModes[i].HashMode < hashModes[j].HashMode
		})
		tbl := uitable.New()
		tbl.AddRow("MODE", "ALGORITHM")
		for _, mode := range hashModes {
			tbl.AddRow(mode.HashMode, mode.Algorithm)
		}
		fmt.Println(tbl)
	},
}

func init() {
	RootCmd.AddCommand(modeCmd)
}
