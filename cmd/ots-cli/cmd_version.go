package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the tool version",
	Run: func(*cobra.Command, []string) {
		fmt.Printf("ots-cli %s\n", version) //nolint:forbidigo
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
