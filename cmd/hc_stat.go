package cmd

import (
	"fmt"
	"sort"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

func getHCStat(f *hashstack.File) {
	path := fmt.Sprintf("/api/hcstat?filename=%s", f.Filename)
	if err := getJSON(path, &f); err != nil {
		writeStdErrAndExit(err.Error())
	}
}

func displayHCStat(f hashstack.File) {
	if f.ID == 0 {
		getHCStat(&f)
	}
	fmt.Printf("Name.............: %s\n", f.Filename)
	fmt.Printf("Size.............: %s\n", humanize.Bytes(uint64(f.Size)))
	fmt.Printf("Last Modified....: %s\n", humanize.Time(time.Unix(f.UpdatedAt, 0)))
	fmt.Println()
}

func displayHCStats() {
	var hcstats []hashstack.File
	if err := getRangeJSON("/api/hcstat", &hcstats); err != nil {
		writeStdErrAndExit(err.Error())
	}
	sort.Slice(hcstats, func(i, j int) bool {
		return hcstats[i].Filename < hcstats[j].Filename
	})
	for _, w := range hcstats {
		displayHCStat(w)
	}
}

var hcstatCmd = &cobra.Command{
	Use:   "hcstats [file_name]",
	Short: "Display a list of all hc_stat files available on the server (-h or --help for subcommands).",
	Long: `
Displays a list of hc_stat files that are stored on the remote server. If file_name is provided, details will be displayed for that specific
hc_stat file.

hc_stat files can be used in jobs. Additional subcommands are available to add and delete files.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			displayHCStat(hashstack.File{Filename: args[0]})
			return
		}
		displayHCStats()
	},
}

var addHCStatCmd = &cobra.Command{
	Use:    "add <file>",
	Short:  "Upload the provided file to the server.",
	Long:   "Upload the provided file to the server.",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("file is required")
		}
		file := args[0]
		uploadFile("/api/hcstat", file)
	},
}

func deleteHCStat(filename string) {
	f := hashstack.File{
		Filename: filename,
	}
	getHCStat(&f)
	path := fmt.Sprintf("/api/hcstat/%d", f.ID)
	if err := deleteHTTP(path); err != nil {
		writeStdErrAndExit(err.Error())
	}
	fmt.Println("hc_stat file deleted successfully")
}

var delHCStatCmd = &cobra.Command{
	Use:    "delete <file_name>",
	Short:  "Delete a file by name from the server.",
	Long:   "Delete a file by name from the server.",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("file_name is required")
		}
		if ok := promptDelete("this file"); !ok {
			writeStdErrAndExit("Not deleting file.")
		}
		deleteHCStat(args[0])
	},
}

func init() {
	hcstatCmd.AddCommand(addHCStatCmd)
	hcstatCmd.AddCommand(delHCStatCmd)
	RootCmd.AddCommand(hcstatCmd)
}
