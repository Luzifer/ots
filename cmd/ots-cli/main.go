package main

import (
	"os"
)

var version = "dev"

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
