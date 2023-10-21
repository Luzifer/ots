package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Luzifer/ots/pkg/customization"
	"github.com/ryanuber/go-glob"
)

var (
	// ErrAttachmentsDisabled signalizes the instance has attachments
	// disabled but the checked secret contains attachments
	ErrAttachmentsDisabled = errors.New("attachments are disabled on this instance")
	// ErrAttachmentsTooLarge signalizes the size of the attached files
	// exceeds the configured maximum size of the given instance
	ErrAttachmentsTooLarge = errors.New("attachment size exceeds allowed size")
	// ErrAttachmentTypeNotAllowed signalizes any file does not match
	// the allowed extensions / mime types
	ErrAttachmentTypeNotAllowed = errors.New("attachment type is not allowed")

	errSettingsNotFound = errors.New("settings not found")
	mimeRegex           = regexp.MustCompile(`^(?:[a-z]+|\*)\/(?:[a-zA-Z0-9.+_-]+|\*)$`)
)

// SanityCheck fetches the instance settings and validates the secret
// against those settings (matching file size, disabled attachments,
// allowed file types, ...)
func SanityCheck(instanceURL string, secret Secret) error {
	cust, err := loadSettings(instanceURL)
	if err != nil {
		if errors.Is(err, errSettingsNotFound) {
			// Sanity check is not possible when the API endpoint is not
			// supported, therefore we ignore this.
			return nil
		}
		return fmt.Errorf("fetching settings: %w", err)
	}

	// Check whether attachments are allowed at all
	if cust.DisableFileAttachment && len(secret.Attachments) > 0 {
		return ErrAttachmentsDisabled
	}

	// Check whether attachments are too large
	var totalAttachmentSize int64
	for _, a := range secret.Attachments {
		totalAttachmentSize += int64(len(a.Content))
	}
	if cust.MaxAttachmentSizeTotal > 0 && totalAttachmentSize > cust.MaxAttachmentSizeTotal {
		return ErrAttachmentsTooLarge
	}

	// Check for allowed types
	if cust.AcceptedFileTypes != "" {
		allowed := strings.Split(cust.AcceptedFileTypes, ",")
		for _, a := range secret.Attachments {
			if !attachmentAllowed(a, allowed) {
				return ErrAttachmentTypeNotAllowed
			}
		}
	}

	return nil
}

func attachmentAllowed(file SecretAttachment, allowed []string) bool {
	for _, a := range allowed {
		switch {
		case mimeRegex.MatchString(a):
			// That's a mime type
			if glob.Glob(a, file.Type) {
				// The mime "glob" matches the file type
				return true
			}

		case a[0] == '.':
			// That's a file extension
			if strings.HasSuffix(file.Name, a) {
				// The filename has the right extension
				return true
			}
		}
	}

	return false
}

func loadSettings(instanceURL string) (c customization.Customize, err error) {
	u, err := url.Parse(instanceURL)
	if err != nil {
		return c, fmt.Errorf("parsing instance URL: %w", err)
	}

	createURL := u.JoinPath(strings.Join([]string{".", "api", "settings"}, "/"))
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, createURL.String(), nil)
	if err != nil {
		return c, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return c, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // possible leaked-fd, lib should not log, potential short-lived leak

	if resp.StatusCode == http.StatusNotFound {
		return c, errSettingsNotFound
	}

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return c, fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
		}
		return c, fmt.Errorf("unexpected HTTP status %d (%s)", resp.StatusCode, respBody)
	}

	if err = json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return c, fmt.Errorf("decoding response: %w", err)
	}

	return c, nil
}
