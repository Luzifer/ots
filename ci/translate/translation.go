package main

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/Luzifer/go_helpers/v2/str"
	"github.com/pkg/errors"
)

type (
	translation     map[string]any
	translationFile struct {
		Reference    translationMapping             `yaml:"reference"`
		Translations map[string]*translationMapping `yaml:"translations"`
	}
	translationMapping struct {
		DeeplLanguage      string      `yaml:"deeplLanguage,omitempty"`
		LanguageKey        string      `yaml:"languageKey,omitempty"`
		Translators        []string    `yaml:"translators"`
		Translations       translation `yaml:"translations"`
		FormalTranslations translation `yaml:"formalTranslations,omitempty"`
	}
)

func (t translation) ToJSON() (string, error) {
	j, err := json.Marshal(t)
	return strings.ReplaceAll(string(j), "'", "\\'"), errors.Wrap(err, "marshalling JSON")
}

func (t translation) GetErrorKeys(ref translation) (missing, extra, wrongType []string) {
	var (
		keys     []string
		keyType  = map[string]reflect.Type{}
		seenKeys []string
	)

	for k, v := range ref {
		keys = append(keys, k)
		keyType[k] = reflect.TypeOf(v)
	}

	for k, v := range t {
		if !str.StringInSlice(k, keys) {
			// Contains extra key, is error
			extra = append(extra, k)
			continue // No further checks for that key
		}

		seenKeys = append(seenKeys, k)
		if kt := reflect.TypeOf(v); keyType[k] != kt {
			// Type mismatches (i.e. string vs []string)
			wrongType = append(wrongType, k)
			continue
		}
	}

	for _, k := range keys {
		if !str.StringInSlice(k, seenKeys) {
			missing = append(missing, k)
		}
	}

	return missing, extra, wrongType
}
