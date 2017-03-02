package cmd

import (
	"fmt"
	"sort"

	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

func getWordlist(f *hashstack.File) {
	path := fmt.Sprintf("/api/wordlists?filename=%s", f.Filename)
	if err := getJSON(path, &f); err != nil {
		writeStdErrAndExit(err.Error())
	}
}

func displayWordlist(f hashstack.File) {
	if f.ID == 0 {
		getWordlist(&f)
	}
	fmt.Printf("Name.............: %s\n", f.Filename)
	fmt.Printf("Words............: %d\n", f.Lines)
	fmt.Printf("Size.............: %s\n", humanize.Bytes(uint64(f.Size)))
	fmt.Printf("Last Modified....: %s\n", humanize.Time(time.Unix(f.UpdatedAt, 0)))
	fmt.Println()
}

func displayWordlists() {
	var wordlists []hashstack.File
	if err := getRangeJSON("/api/wordlists", &wordlists); err != nil {
		writeStdErrAndExit(err.Error())
	}
	sort.Slice(wordlists, func(i, j int) bool {
		return wordlists[i].Filename < wordlists[j].Filename
	})
	for _, w := range wordlists {
		displayWordlist(w)
	}
}

var wordlistCmd = &cobra.Command{
	Use:   "wordlists [file_name]",
	Short: "Display a list of all wordlists available on the server (-h or --help for subcommands)",
	Long: `
Displays a list of wordlists that are stored on the remote server. If file_name is provided, details will be displayed for that specific
wordlist.

Wordlists can be used in jobs. Additional subcommands are available to add and delete files.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			displayWordlist(hashstack.File{Filename: args[0]})
			return
		}
		displayWordlists()
	},
}

var addWordlistCmd = &cobra.Command{
	Use:    "add <file>",
	Short:  "Upload the provided file to the server",
	Long:   "Upload the provided file to the server",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("file is required")
		}
		file := args[0]
		uploadFile("/api/wordlists", file)
	},
}

func deleteWordlist(filename string) {
	f := hashstack.File{
		Filename: filename,
	}
	getWordlist(&f)
	path := fmt.Sprintf("/api/wordlists/%d", f.ID)
	if err := deleteHTTP(path); err != nil {
		writeStdErrAndExit(err.Error())
	}
	fmt.Println("wordlist deleted successfully")
}

var delWordlistCmd = &cobra.Command{
	Use:    "delete <file_name>",
	Short:  "Delete a file by name from the server",
	Long:   "Delete a file by name from the server",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("file_name is required")
		}
		deleteWordlist(args[0])
	},
}

func init() {
	wordlistCmd.AddCommand(addWordlistCmd)
	wordlistCmd.AddCommand(delWordlistCmd)
	RootCmd.AddCommand(wordlistCmd)
}
