package main

import (
	"fmt"
	"os"

	"gogen/internal/cli"
)

var (
	version = "1.0.0"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	rootCmd := cli.NewRootCommand()

	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
