package main

import (
	"github.com/spf13/cobra/doc"
	"github.com/stricture/hashstack-cli/cmd"
)

func main() {
	cmd.RootCmd.DisableAutoGenTag = true
	doc.GenMarkdownTree(cmd.RootCmd, "./markdown")
	header := &doc.GenManHeader{
		Title:   "hashstack",
		Section: "1",
		Source:  "Terahash",
	}
	doc.GenManTree(cmd.RootCmd, header, "./man")
}
