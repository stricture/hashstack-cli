package cmd

import (
	"fmt"
	"sort"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

func getRule(f *hashstack.File) {
	path := fmt.Sprintf("/api/rules?filename=%s", f.Filename)
	if err := getJSON(path, &f); err != nil {
		writeStdErrAndExit(err.Error())
	}
}

func displayRule(f hashstack.File) {
	if f.ID == 0 {
		getRule(&f)
	}
	fmt.Printf("Name.............: %s\n", f.Filename)
	fmt.Printf("Rules............: %d\n", f.Lines)
	fmt.Printf("Size.............: %s\n", humanize.Bytes(uint64(f.Size)))
	fmt.Printf("Last Modified....: %s\n", humanize.Time(time.Unix(f.UpdatedAt, 0)))
	fmt.Println()
}

func displayRules() {
	var rules []hashstack.File
	if err := getRangeJSON("/api/rules", &rules); err != nil {
		writeStdErrAndExit(err.Error())
	}
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Filename < rules[j].Filename
	})
	for _, w := range rules {
		displayRule(w)
	}
}

var ruleCmd = &cobra.Command{
	Use:   "rules [file_name]",
	Short: "Display a list of all rule files available on the server (-h or --help for subcommands).",
	Long: `
Displays a list of rule files that are stored on the remote server. If file_name is provided, details will be displayed for that specific
rule file.

Rule files can be used in jobs. Additional subcommands are available to add and delete files.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			displayRule(hashstack.File{Filename: args[0]})
			return
		}
		displayRules()
	},
}

var addRuleCmd = &cobra.Command{
	Use:    "add <file>",
	Short:  "Upload the provided file to the server.",
	Long:   "Upload the provided file to the server.",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("file is required.")
		}
		file := args[0]
		uploadFile("/api/rules", file)
	},
}

func deleteRule(filename string) {
	f := hashstack.File{
		Filename: filename,
	}
	getRule(&f)
	path := fmt.Sprintf("/api/rules/%d", f.ID)
	if err := deleteHTTP(path); err != nil {
		writeStdErrAndExit(err.Error())
	}
	fmt.Println("The rule was deleted successfully.")
}

var delRuleCmd = &cobra.Command{
	Use:    "delete <file_name>",
	Short:  "Delete a file by name from the server.",
	Long:   "Delete a file by name from the server.",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("file_name is required.")
		}
		if ok := promptDelete("this file"); !ok {
			writeStdErrAndExit("Not deleting file.")
		}
		deleteRule(args[0])
	},
}

func init() {
	ruleCmd.AddCommand(addRuleCmd)
	ruleCmd.AddCommand(delRuleCmd)
	RootCmd.AddCommand(ruleCmd)
}
