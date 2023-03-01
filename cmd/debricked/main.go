package main

import (
	"os"

	"github.com/debricked/cli/pkg/cmd/root"
)

var version string // Set at compile time

func main() {
	if err := root.NewRootCmd(version).Execute(); err != nil {
		os.Exit(1)
	}
}
