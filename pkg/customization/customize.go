// Package customization contains the structure for the customization
// file to configure the OTS web- and command-line interface
package customization

import (
	"encoding/json"
	"io/fs"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Frontend has a max attachment size of 64MiB as the base64 encoding
// will break afterwards. Therefore we use a maximum secret size of
// 65MiB and increase it by double base64 encoding:
//
// 65 MiB * 16/9 (twice 4/3 base64 size increase)
const defaultMaxSecretSize = 65 * 1024 * 1024 * (16 / 9) // = 115.6MiB

type (
	// Customize holds the structure of the customization file
	Customize struct {
		AppIcon              string `json:"appIcon,omitempty" yaml:"appIcon"`
		AppTitle             string `json:"appTitle,omitempty" yaml:"appTitle"`
		DisableAppTitle      bool   `json:"disableAppTitle,omitempty" yaml:"disableAppTitle"`
		DisablePoweredBy     bool   `json:"disablePoweredBy,omitempty" yaml:"disablePoweredBy"`
		DisableQRSupport     bool   `json:"disableQRSupport,omitempty" yaml:"disableQRSupport"`
		DisableThemeSwitcher bool   `json:"disableThemeSwitcher,omitempty" yaml:"disableThemeSwitcher"`

		DisableExpiryOverride bool    `json:"disableExpiryOverride,omitempty" yaml:"disableExpiryOverride"`
		ExpiryChoices         []int64 `json:"expiryChoices,omitempty" yaml:"expiryChoices"`

		AcceptedFileTypes      string `json:"acceptedFileTypes" yaml:"acceptedFileTypes"`
		DisableFileAttachment  bool   `json:"disableFileAttachment" yaml:"disableFileAttachment"`
		MaxAttachmentSizeTotal int64  `json:"maxAttachmentSizeTotal" yaml:"maxAttachmentSizeTotal"`

		MaxSecretSize         int64    `json:"-" yaml:"maxSecretSize"`
		MetricsAllowedSubnets []string `json:"-" yaml:"metricsAllowedSubnets"`
		OverlayFSPath         string   `json:"-" yaml:"overlayFSPath"`
		UseFormalLanguage     bool     `json:"-" yaml:"useFormalLanguage"`
	}
)

// Load retrieves the Customization file from filesystem
func Load(filename string) (cust Customize, err error) {
	if filename == "" {
		// None given, take a shortcut
		cust.applyFixes()
		return cust, nil
	}

	cf, err := os.Open(filename) //#nosec:G304 // Loading a custom file is the intention here
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logrus.Warn("customize file given but not found")
			return cust, nil
		}
		return cust, errors.Wrap(err, "opening customize file")
	}
	defer func() {
		if err := cf.Close(); err != nil {
			logrus.WithError(err).Error("closing customize file (leaked fd)")
		}
	}()

	if err = yaml.NewDecoder(cf).Decode(&cust); err != nil {
		return cust, errors.Wrap(err, "decoding customize file")
	}

	cust.applyFixes()

	return cust, nil
}

// ToJSON is a templating helper which returns the customization
// serialized as JSON in a string
func (c Customize) ToJSON() (string, error) {
	j, err := json.Marshal(c)
	return string(j), errors.Wrap(err, "marshalling JSON")
}

func (c *Customize) applyFixes() {
	if len(c.AppTitle) == 0 {
		c.AppTitle = "OTS - One Time Secrets"
	}

	if c.MaxSecretSize == 0 {
		c.MaxSecretSize = defaultMaxSecretSize
	}
}
