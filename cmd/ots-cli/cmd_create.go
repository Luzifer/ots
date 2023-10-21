package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/Luzifer/ots/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type (
	authRoundTripper struct {
		http.RoundTripper

		headers    http.Header
		user, pass string
	}
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
	createCmd.Flags().StringSliceP("header", "H", nil, "Headers to include in the request (i.e. 'Authorization: Token ...')")
	createCmd.Flags().String("instance", "https://ots.fyi/", "Instance to create the secret with")
	createCmd.Flags().StringSliceP("file", "f", nil, "File(s) to attach to the secret")
	createCmd.Flags().String("secret-from", "-", `File to read the secret content from ("-" for STDIN)`)
	createCmd.Flags().StringP("user", "u", "", "Username / Password for basic auth, specified as 'user:pass'")
	rootCmd.AddCommand(createCmd)
}

func createRunE(cmd *cobra.Command, _ []string) (err error) {
	var secret client.Secret

	if client.HTTPClient, err = constructHTTPClient(cmd); err != nil {
		return fmt.Errorf("constructing authorized HTTP client: %w", err)
	}

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
			Name:    path.Base(f),
			Type:    mime.TypeByExtension(path.Ext(f)),
			Content: content,
		})
	}

	// Get flags for creation
	logrus.Info("creating the secret...")
	instanceURL, err := cmd.Flags().GetString("instance")
	if err != nil {
		return fmt.Errorf("getting instance flag: %w", err)
	}

	expire, err := cmd.Flags().GetDuration("expire")
	if err != nil {
		return fmt.Errorf("getting expire flag: %w", err)
	}

	// Execute sanity checks
	if err = client.SanityCheck(instanceURL, secret); err != nil {
		return fmt.Errorf("sanity checking secret: %w", err)
	}

	// Create the secret
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

func constructHTTPClient(cmd *cobra.Command) (*http.Client, error) {
	basic, _ := cmd.Flags().GetString("user")
	headers, _ := cmd.Flags().GetStringSlice("header")

	if basic == "" && headers == nil {
		// No authorization needed
		return http.DefaultClient, nil
	}

	t := authRoundTripper{RoundTripper: http.DefaultTransport, headers: http.Header{}}

	// Set basic auth if available
	user, pass, ok := strings.Cut(basic, ":")
	if ok {
		t.user = user
		t.pass = pass
	}

	// Parse and set headers if available
	for _, hdr := range headers {
		key, value, ok := strings.Cut(hdr, ":")
		if !ok {
			logrus.WithField("header", hdr).Warn("invalid header format, skipping")
			continue
		}
		t.headers.Add(key, strings.TrimSpace(value))
	}

	return &http.Client{Transport: t}, nil
}

func (a authRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if a.user != "" {
		r.SetBasicAuth(a.user, a.pass)
	}

	for key, values := range a.headers {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}

	resp, err := a.RoundTripper.RoundTrip(r)
	if err != nil {
		return nil, fmt.Errorf("executing round-trip: %w", err)
	}
	return resp, nil
}
