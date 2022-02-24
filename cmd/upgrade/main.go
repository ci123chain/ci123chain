package main

import (
	"fmt"
	"github.com/ci123chain/ci123chain/cmd/upgrade/upgrade"
	"os"
)

func main() {
	if err := Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

// Run is the main loop, but returns an error
func Run(args []string) error {
	cfg, err := upgrade.GetConfigFromEnv()
	if err != nil {
		return err
	}
	return upgrade.LaunchProcess(cfg, args, os.Stdout, os.Stderr)
}
