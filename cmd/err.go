package cmd

import (
	"fmt"
	"os"
)

func writeStdErrAndExit(msg string) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", msg))
	if !flDebug {
		fmt.Fprintf(os.Stderr, "\nRunning with --debug will show additional context for this error.\n")
	}
	fmt.Fprintf(os.Stderr, "\nEmail support@terahash.com for more information. Please include --debug output!\n")
	os.Exit(1)
}
