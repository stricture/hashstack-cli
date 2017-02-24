package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

var fileCmd = &cobra.Command{
	Use:    "files",
	Short:  "Subcommands can be used to interact with hashstack files",
	Long:   "Subcommands can be used to interact with hashstack files",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("try hashstack-cli files -h")
	},
}

var listFileCmd = &cobra.Command{
	Use:    "list <type>",
	Short:  "Prints a list of files on the remote server",
	Long:   "Prints a list of files on the remote server",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("<type> is required")
		}
		var files []hashstack.File
		switch args[0] {
		case "wordlists":
			if err := getRangeJSON("/api/wordlists", &files); err != nil {
				writeStdErrAndExit(err.Error())
			}
		case "rules":
			if err := getRangeJSON("/api/rules", &files); err != nil {
				writeStdErrAndExit(err.Error())
			}
		default:
			writeStdErrAndExit("<type> is not valid")
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].Filename < files[j].Filename
		})
		tbl := uitable.New()
		tbl.AddRow("FILENAME", "LINES", "MTIME")
		for _, file := range files {
			tbl.AddRow(file.Filename, file.Lines, time.Unix(file.Mtime, 0).String())
		}
		fmt.Println(tbl)
	},
}

func init() {
	fileCmd.AddCommand(listFileCmd)
	RootCmd.AddCommand(fileCmd)
}
