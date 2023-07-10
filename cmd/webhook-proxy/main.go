package main

import (
	"os"

	"github.com/snapp-incubator/jira-element-proxy/internal/webhook-proxy/cmd"
)

const (
	exitFailure = 1
)

func main() {
	root := cmd.NewRootCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(exitFailure)
		}
	}
}
