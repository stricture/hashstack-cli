package main

import (
	"github.com/spf13/cobra/doc"
	"github.com/stricture/hashstack-cli/cmd"
)

func main() {
	doc.GenMarkdownTree(cmd.RootCmd, "./markdown")
	header := &doc.GenManHeader{
		Title:   "hashstack",
		Section: "1",
	}
	header.Source = "Terahash"
	doc.GenManTree(cmd.RootCmd, header, "./man")
}
