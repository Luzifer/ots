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
		AppIcon               string  `json:"appIcon,omitempty" yaml:"appIcon"`
		AppTitle              string  `json:"appTitle,omitempty" yaml:"appTitle"`
		DisableAppTitle       bool    `json:"disableAppTitle,omitempty" yaml:"disableAppTitle"`
		DisableExpiryOverride bool    `json:"disableExpiryOverride,omitempty" yaml:"disableExpiryOverride"`
		DisablePoweredBy      bool    `json:"disablePoweredBy,omitempty" yaml:"disablePoweredBy"`
		DisableQRSupport      bool    `json:"disableQRSupport,omitempty" yaml:"disableQRSupport"`
		DisableThemeSwitcher  bool    `json:"disableThemeSwitcher,omitempty" yaml:"disableThemeSwitcher"`
		ExpiryChoices         []int64 `json:"expiryChoices,omitempty" yaml:"expiryChoices"`
		OverlayFSPath         string  `json:"-" yaml:"overlayFSPath"`
		UseFormalLanguage     bool    `json:"-" yaml:"useFormalLanguage"`
	}
)

func loadCustomize(filename string) (customize, error) {
	if filename == "" {
		// None given, take a shortcut
		return customize{}, nil
	}

	var cust customize

	cf, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logrus.Warn("customize file given but not found")
			return cust, nil
		}
		return cust, errors.Wrap(err, "opening customize file")
	}
	defer cf.Close()

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
