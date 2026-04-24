package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Luzifer/ots/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const storeFileMode = 0o600 // We assume the attached file to be a secret

var fetchCmd = &cobra.Command{
	Use:   "fetch <url>",
	Short: "Retrieves a secret from the instance by its URL",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	RunE:  fetchRunE,
}

func init() {
	fetchCmd.Flags().String("file-dir", ".", "Where to put files attached to the secret")
	rootCmd.AddCommand(fetchCmd)
}

func assembleDownloadFileName(dir, filename string, iteration int) string {
	fileNameFragments := strings.SplitN(filepath.Base(filename), ".", 2) //nolint:mnd

	switch {
	case iteration == 0 && len(fileNameFragments) == 1:
		// We are in initial iteration and have no file extension
		return filepath.Join(dir, fileNameFragments[0])

	case iteration == 0:
		// Initial iteration, extension is present
		return filepath.Join(dir, strings.Join(fileNameFragments, "."))

	case len(fileNameFragments) == 1:
		// Later iteration and no extension
		return filepath.Join(dir, fmt.Sprintf("%s (%d)", fileNameFragments[0], iteration))

	default:
		// Later iteration, extension is present
		return filepath.Join(dir, fmt.Sprintf("%s (%d).%s", fileNameFragments[0], iteration, fileNameFragments[1]))
	}
}

func checkDirWritable(dir string) error {
	tmpFile := filepath.Join(dir, ".ots-cli.tmp")
	if err := os.WriteFile(tmpFile, []byte(""), storeFileMode); err != nil {
		return fmt.Errorf("writing tmp-file: %w", err)
	}
	defer os.Remove(tmpFile) //nolint:errcheck // We don't really care

	return nil
}

func fetchRunE(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	fileDir, err := cmd.Flags().GetString("file-dir")
	if err != nil {
		return fmt.Errorf("getting file-dir parameter: %w", err)
	}

	// First lets check whether we potentially can write files
	if err := checkDirWritable(fileDir); err != nil {
		return fmt.Errorf("checking for directory write: %w", err)
	}

	logrus.Info("fetching secret...")
	secret, err := client.Fetch(args[0])
	if err != nil {
		return fmt.Errorf("fetching secret")
	}

	for _, f := range secret.Attachments {
		logrus.WithField("file", f.Name).Info("storing file...")
		if err = storeAttachment(fileDir, f); err != nil {
			return fmt.Errorf("saving file to disk: %w", err)
		}
	}

	fmt.Println(secret.Secret) //nolint:forbidigo // Output intended for STDOUT

	return nil
}

func storeAttachment(dir string, f client.SecretAttachment) (err error) {
	// First lets find a free file name to save the file as
	var (
		i         int
		storeName string
	)

	if slices.Contains([]string{"", ".", "/", `\`}, filepath.Base(f.Name)) {
		// These "filenames" makes no sense and could cause trouble when storing
		return fmt.Errorf("invalid attachment name %q", f.Name)
	}

	for {
		storeName = assembleDownloadFileName(dir, f.Name, i)
		if _, err = os.Stat(storeName); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				break
			}

			return fmt.Errorf("getting file stat: %w", err)
		}

		// No luck, file is taken, next round
		i++
	}

	// So we finally found a filename we can use
	if err := os.WriteFile(storeName, f.Content, storeFileMode); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
