package cmd

import (
	"github.com/segmentio/go-prompt"
)

func promptDelete(thing string) bool {
	return prompt.Confirm("Are you sure you want to delete %s? [yY/nN]", thing)
}
