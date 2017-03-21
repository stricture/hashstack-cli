package main

import (
	"github.com/spf13/cobra/doc"
	"github.com/stricture/hashstack-cli/cmd"
)

func main() {
	doc.GenMarkdownTree(cmd.RootCmd, "./")
	header := &doc.GenManHeader{
		Title:   "Hashstack-cli",
		Section: "1",
	}
	doc.GenManTree(cmd.RootCmd, header, "./")
}
