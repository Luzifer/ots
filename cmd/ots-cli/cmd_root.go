package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short:             "Utility to interact with encrypted secrets in an OTS instance",
	PersistentPreRunE: rootPersistentPreRunE,
}

func init() {
	rootCmd.PersistentFlags().String("log-level", "info", "Level to use for logging (trace, debug, info, warn, error, fatal)")
}

func rootPersistentPreRunE(cmd *cobra.Command, _ []string) error {
	sll, err := cmd.Flags().GetString("log-level")
	if err != nil {
		return fmt.Errorf("getting log-level: %w", err)
	}

	ll, err := logrus.ParseLevel(sll)
	if err != nil {
		return fmt.Errorf("parsing log-level: %w", err)
	}
	logrus.SetLevel(ll)

	return nil
}
