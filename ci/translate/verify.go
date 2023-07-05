package main

import (
	"reflect"
	"regexp"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/go_helpers/v2/str"
)

var langKeyFormat = regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`)

func verify(tf translationFile) error {
	var (
		err     error
		keys    []string
		keyType = map[string]reflect.Type{}
	)

	for k, v := range tf.Reference.Translations {
		keys = append(keys, k)
		keyType[k] = reflect.TypeOf(v)
	}

	if !langKeyFormat.MatchString(tf.Reference.LanguageKey) {
		return errors.New("reference contains invalid languageKey")
	}

	if len(keys) == 0 {
		return errors.New("reference does not contain translations")
	}

	logrus.Infof("found %d translation keys in reference", len(keys))

	if tf.Reference.FormalTranslations != nil {
		if verifyTranslationKeys(logrus.NewEntry(logrus.StandardLogger()), tf.Reference.FormalTranslations, keys, keyType, false); err != nil {
			return errors.New("reference contains error in formalTranslations")
		}
	}

	var hadErrors bool
	for lk, tm := range tf.Translations {
		logger := logrus.WithField("lang", lk)
		logger.Info("validating language...")

		if !langKeyFormat.MatchString(lk) {
			hadErrors = true
			logger.Error("language key is invalid")
		}

		if tm.DeeplLanguage == "" {
			logger.Info("no deeplLanguage is set")
		}

		hadErrors = hadErrors || verifyTranslationKeys(logger, tm.Translations, keys, keyType, true)
		hadErrors = hadErrors || verifyTranslationKeys(logger, tm.FormalTranslations, keys, keyType, false)
	}

	if hadErrors {
		return errors.New("translation file has errors")
	}
	return nil
}

func verifyTranslationKeys(logger *logrus.Entry, t translation, keys []string, keyType map[string]reflect.Type, warnMissing bool) (hadErrors bool) {
	var seenKeys []string

	for k, v := range t {
		keyLogger := logger.WithField("translation_key", k)
		if !str.StringInSlice(k, keys) {
			// Contains extra key, is error
			hadErrors = true
			keyLogger.Error("extra key found")
			continue // No further checks for that key
		}

		seenKeys = append(seenKeys, k)
		if kt := reflect.TypeOf(v); keyType[k] != kt {
			// Type mismatches (i.e. string vs []string)
			hadErrors = true
			keyLogger.Errorf("key has invalid type %s != %s", kt, keyType[k])
		}
	}

	for _, k := range keys {
		if warnMissing && !str.StringInSlice(k, seenKeys) {
			logger.WithField("translation_key", k).Warn("missing translation")
		}
	}

	return hadErrors
}
