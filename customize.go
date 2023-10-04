package main

import (
	"encoding/json"
	"io/fs"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type (
	customize struct {
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

		OverlayFSPath     string `json:"-" yaml:"overlayFSPath"`
		UseFormalLanguage bool   `json:"-" yaml:"useFormalLanguage"`
	}
)

func loadCustomize(filename string) (cust customize, err error) {
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

func (c customize) ToJSON() (string, error) {
	j, err := json.Marshal(c)
	return string(j), errors.Wrap(err, "marshalling JSON")
}

func (c *customize) applyFixes() {
	if len(c.AppTitle) == 0 {
		c.AppTitle = "OTS - One Time Secrets"
	}
}
