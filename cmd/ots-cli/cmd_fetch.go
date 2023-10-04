package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/Luzifer/ots/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const storeFileMode = 0o600 // We assume the attached file to be a secret

var fetchCmd = &cobra.Command{
	Use:   "fetch url",
	Short: "Retrieves a secret from the instance by its URL",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	RunE:  fetchRunE,
}

func init() {
	fetchCmd.Flags().String("file-dir", ".", "Where to put files attached to the secret")
	rootCmd.AddCommand(fetchCmd)
}

func checkDirWritable(dir string) error {
	tmpFile := path.Join(dir, ".ots-cli.tmp")
	if err := os.WriteFile(tmpFile, []byte(""), storeFileMode); err != nil {
		return fmt.Errorf("writing tmp-file: %w", err)
	}
	defer os.Remove(tmpFile) //nolint:errcheck // We don't really care

	return nil
}

func fetchRunE(cmd *cobra.Command, args []string) error {
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

func storeAttachment(dir string, f client.SecretAttachment) error {
	// First lets find a free file name to save the file as
	var (
		fileNameFragments = strings.SplitN(f.Name, ".", 2) //nolint:gomnd
		i                 int
		storeName         = path.Join(dir, f.Name)
		storeNameTpl      string
	)

	if len(fileNameFragments) == 1 {
		storeNameTpl = fmt.Sprintf("%s (%%d)", fileNameFragments[0])
	} else {
		storeNameTpl = fmt.Sprintf("%s (%%d).%s", fileNameFragments[0], fileNameFragments[1])
	}

	for _, err := os.Stat(storeName); !errors.Is(err, fs.ErrNotExist); _, err = os.Stat(storeName) {
		i++
		storeName = fmt.Sprintf(storeNameTpl, i)
	}

	// So we finally found a filename we can use
	if err := os.WriteFile(storeName, f.Content, storeFileMode); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
