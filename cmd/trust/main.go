package main

import (
	"os"

	"github.com/Darren-Certlocker/trust-cli/internal/cli"
)

func main() {
	if err := cli.New().Execute(); err != nil {
		os.Exit(1)
	}
}
