package cmd

import (
	"fmt"
	"os"
)

func writeStdErrAndExit(msg string) {
	fmt.Fprintf(os.Stderr, msg)
	os.Exit(1)
}
