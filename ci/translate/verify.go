package main

import (
	"regexp"
	"sort"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var langKeyFormat = regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`)

func verify(tf translationFile) error {
	var err error

	if !langKeyFormat.MatchString(tf.Reference.LanguageKey) {
		return errors.New("reference contains invalid languageKey")
	}

	if len(tf.Reference.Translations) == 0 {
		return errors.New("reference does not contain translations")
	}

	logrus.Infof("found %d translation keys in reference", len(tf.Reference.Translations))

	if tf.Reference.FormalTranslations != nil {
		if verifyTranslationKeys(
			logrus.NewEntry(logrus.StandardLogger()),
			tf.Reference.FormalTranslations,
			tf.Reference.Translations,
			false,
		); err != nil {
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

		hadErrors = hadErrors || verifyTranslationKeys(
			logger,
			tm.Translations,
			tf.Reference.Translations,
			true,
		)
		hadErrors = hadErrors || verifyTranslationKeys(
			logger,
			tm.FormalTranslations,
			tf.Reference.Translations,
			false,
		)
	}

	if hadErrors {
		return errors.New("translation file has errors")
	}
	return nil
}

//revive:disable-next-line:flag-parameter
func verifyTranslationKeys(logger *logrus.Entry, t, ref translation, warnMissing bool) (hadErrors bool) {
	missing, extra, wrongType := t.GetErrorKeys(ref)

	sort.Strings(extra)
	sort.Strings(missing)
	sort.Strings(wrongType)

	for _, k := range extra {
		logger.WithField("translation_key", k).Error("extra key found")
	}

	for _, k := range wrongType {
		logger.WithField("translation_key", k).Error("key has invalid type")
	}

	if warnMissing {
		for _, k := range missing {
			logger.WithField("translation_key", k).Warn("missing translation")
		}
	}

	return len(extra)+len(wrongType) > 0
}
