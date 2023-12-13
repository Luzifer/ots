package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the tool version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ots-cli %s\n", version) //nolint:forbidigo
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
