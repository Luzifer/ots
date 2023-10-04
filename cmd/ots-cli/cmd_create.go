package main

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path"

	"github.com/Luzifer/ots/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create [-f file]... [--instance url] [--secret-from file]",
	Short:   "Create a new encrypted secret in the given OTS instance",
	Long:    "",
	Example: `echo "I'm a very secret secret" | ots-cli create`,
	Args:    cobra.NoArgs,
	RunE:    createRunE,
}

func init() {
	createCmd.Flags().Duration("expire", 0, "When to expire the secret (0 to use server-default)")
	createCmd.Flags().String("instance", "https://ots.fyi/", "Instance to create the secret with")
	createCmd.Flags().StringSliceP("file", "f", nil, "File(s) to attach to the secret")
	createCmd.Flags().String("secret-from", "-", `File to read the secret content from ("-" for STDIN)`)
	rootCmd.AddCommand(createCmd)
}

func createRunE(cmd *cobra.Command, _ []string) error {
	var secret client.Secret

	// Read the secret content
	logrus.Info("reading secret content...")
	secretSourceName, err := cmd.Flags().GetString("secret-from")
	if err != nil {
		return fmt.Errorf("getting secret-from flag: %w", err)
	}

	var secretSource io.Reader
	if secretSourceName == "-" {
		secretSource = os.Stdin
	} else {
		f, err := os.Open(secretSourceName) //#nosec:G304 // Opening user specified file is intended
		if err != nil {
			return fmt.Errorf("opening secret-from file: %w", err)
		}
		defer f.Close() //nolint:errcheck // The file will be force-closed by program exit
		secretSource = f
	}

	secretContent, err := io.ReadAll(secretSource)
	if err != nil {
		return fmt.Errorf("reading secret content: %w", err)
	}
	secret.Secret = string(secretContent)

	// Attach any file given
	files, err := cmd.Flags().GetStringSlice("file")
	if err != nil {
		return fmt.Errorf("getting file flag: %w", err)
	}
	for _, f := range files {
		logrus.WithField("file", f).Info("attaching file...")
		content, err := os.ReadFile(f) //#nosec:G304 // Opening user specified file is intended
		if err != nil {
			return fmt.Errorf("reading attachment %q: %w", f, err)
		}

		secret.Attachments = append(secret.Attachments, client.SecretAttachment{
			Name:    f,
			Type:    mime.TypeByExtension(path.Ext(f)),
			Content: content,
		})
	}

	// Create the secret
	logrus.Info("creating the secret...")
	instanceURL, err := cmd.Flags().GetString("instance")
	if err != nil {
		return fmt.Errorf("getting instance flag: %w", err)
	}

	expire, err := cmd.Flags().GetDuration("expire")
	if err != nil {
		return fmt.Errorf("getting expire flag: %w", err)
	}

	secretURL, expiresAt, err := client.Create(instanceURL, secret, expire)
	if err != nil {
		return fmt.Errorf("creating secret: %w", err)
	}

	// Tell them where to find the secret
	if expiresAt.IsZero() {
		logrus.Info("secret created, see URL below")
	} else {
		logrus.WithField("expires-at", expiresAt).Info("secret created, see URL below")
	}
	fmt.Println(secretURL) //nolint:forbidigo // Output intended for STDOUT

	return nil
}
