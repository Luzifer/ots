package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/Luzifer/rconfig/v2"
)

const deeplRequestTimeout = 10 * time.Second

var (
	cfg = struct {
		AutoTranslate    bool   `flag:"auto-translate" default:"false" description:"Enable auto-translation through DeepL"`
		DeeplAPIEndpoint string `flag:"deepl-api-endpoint" default:"https://api-free.deepl.com/v2/translate" description:"DeepL API endpoint to request translations from"`
		DeeplAPIKey      string `flag:"deepl-api-key" default:"" description:"API key for the DeepL API"`
		IssueFile        string `flag:"issue-file" default:"../../translate-issue.md" description:"Where to create the translate issue"`
		OutputFile       string `flag:"output-file,o" default:"../../src/langs/langs.js" description:"Where to put rendered translations"`
		Template         string `flag:"template" default:"../../src/langs/langs.tpl.js" description:"Template to load for translation JS file"`
		TranslationFile  string `flag:"translation-file,t" default:"../../i18n.yaml" description:"File to use for translations"`
		LogLevel         string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		Verify           bool   `flag:"verify" default:"true" description:"Run verification against translation file"`
		VersionAndExit   bool   `flag:"version" default:"false" description:"Prints current version and exits"`
		WriteIssueFile   bool   `flag:"write-issue-file" default:"false" description:"Generates an issue body for missing translations"`
	}{}

	version = "dev"
)

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return errors.Wrap(err, "parsing cli options")
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return errors.Wrap(err, "parsing log-level")
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	if cfg.VersionAndExit {
		logrus.WithField("version", version).Info("translate")
		os.Exit(0)
	}

	logrus.Info("loading translations...")

	tf, err := loadTranslationFile()
	if err != nil {
		logrus.WithError(err).Fatal("loading translation file")
	}

	if cfg.AutoTranslate {
		logrus.Info("auto-translating new strings...")

		if err = autoTranslate(&tf); err != nil {
			logrus.WithError(err).Fatal("adding missing translations")
		}
	}

	if cfg.Verify {
		logrus.Info("verify translation file...")

		if err = verify(tf); err != nil {
			logrus.WithError(err).Fatal("verifying translations")
		}
	}

	logrus.Info("saving translation file...")

	if err = saveTranslationFile(tf); err != nil {
		logrus.WithError(err).Fatal("saving translation file")
	}

	logrus.Info("updating JS embedded translations...")

	// Copy reference for rendering
	tf.Translations[tf.Reference.LanguageKey] = &tf.Reference

	if err = renderJSFile(tf); err != nil {
		logrus.WithError(err).Fatal("rendering JS output")
	}

	if cfg.WriteIssueFile {
		logrus.Info("writing issue template...")
		if err = generateIssue(tf); err != nil {
			logrus.WithError(err).Fatal("generating issue template")
		}
	}
}

func autoTranslate(tf *translationFile) error {
	if cfg.DeeplAPIKey == "" {
		logrus.Warn("missing DeepL API key, skipping translation of new strings")
		return nil
	}

	// Collect keys to translate
	var keys []string
	for key := range tf.Reference.Translations {
		keys = append(keys, key)
	}

	for lang := range tf.Translations {
		if tf.Translations[lang].DeeplLanguage == "" {
			logrus.WithField("lang", lang).Warn("missing DeepL language, skipping")
			continue
		}

		for _, key := range keys {
			if err := autoTranslateKeyForLang(tf, lang, key); err != nil {
				return errors.Wrapf(err, "translating %s:%s", lang, key)
			}
		}
	}

	return nil
}

func autoTranslateKeyForLang(tf *translationFile, lang, key string) (err error) {
	if tf.Translations[lang].Translations[key] != nil {
		// There is something, we assume that's fine, might miss out newly
		// added strings in a slice - we care about that when we need to
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"lang": lang,
		"key":  key,
	}).Info("fetching translation...")

	if tf.Translations[lang].Translations == nil {
		tf.Translations[lang].Translations = make(map[string]any)
	}

	switch typedSrc := tf.Reference.Translations[key].(type) {
	case string:
		if tf.Translations[lang].Translations[key], err = fetchTranslation(
			tf.Reference.DeeplLanguage,
			tf.Translations[lang].DeeplLanguage,
			typedSrc,
		); err != nil {
			return errors.Wrapf(err, "translating %s:%s", lang, key)
		}

	case []any:
		var ts []string
		for _, str := range typedSrc {
			tStr, err := fetchTranslation(
				tf.Reference.DeeplLanguage,
				tf.Translations[lang].DeeplLanguage,
				str.(string),
			)
			if err != nil {
				return errors.Wrapf(err, "translating %s:%s", lang, key)
			}
			ts = append(ts, tStr)
		}
		tf.Translations[lang].Translations[key] = ts

	default:
		return errors.Errorf("unexpected translation type %T", tf.Reference.Translations[key])
	}

	return nil
}

func fetchTranslation(srcLang, destLang, text string) (string, error) {
	params := url.Values{}
	params.Set("text", text)
	params.Set("source_lang", strings.ToUpper(srcLang))
	params.Set("target_lang", strings.ToUpper(destLang))
	params.Set("tag_handling", "html")

	ctx, cancel := context.WithTimeout(context.Background(), deeplRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.DeeplAPIEndpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return "", errors.Wrap(err, "creating request")
	}
	req.Header.Set("Authorization", strings.Join([]string{"DeepL-Auth-Key", cfg.DeeplAPIKey}, " "))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "executing request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.WithError(err).Error("closing response body (leaked fd)")
		}
	}()

	var payload struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", errors.Wrap(err, "decoding DeepL response")
	}

	if l := len(payload.Translations); l != 1 {
		return "", errors.Errorf("unexpected number of translations: %d", l)
	}

	return payload.Translations[0].Text, nil
}

func loadTranslationFile() (translationFile, error) {
	var tf translationFile
	f, err := os.Open(cfg.TranslationFile)
	if err != nil {
		return tf, errors.Wrap(err, "opening translation file")
	}
	defer f.Close() //nolint:errcheck // Short-lived fd-leak

	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)

	return tf, errors.Wrap(decoder.Decode(&tf), "decoding translation file")
}

func renderJSFile(tf translationFile) error {
	jsTemplate, err := os.ReadFile(cfg.Template)
	if err != nil {
		return errors.Wrap(err, "reading template file")
	}

	tpl, err := template.New("js").Funcs(sprig.FuncMap()).Parse(string(jsTemplate))
	if err != nil {
		return errors.Wrap(err, "parsing template")
	}

	f, err := os.Create(cfg.OutputFile + ".tmp")
	if err != nil {
		return errors.Wrap(err, "creating tempfile")
	}

	if err = tpl.Execute(f, tf); err != nil {
		f.Close() //nolint:errcheck,gosec,revive // Short-lived fd-leak
		return errors.Wrap(err, "rendering js template")
	}

	f.Close() //nolint:errcheck,gosec,revive // Short-lived fd-leak
	return errors.Wrap(os.Rename(cfg.OutputFile+".tmp", cfg.OutputFile), "moving file in place")
}

func saveTranslationFile(tf translationFile) error {
	f, err := os.Create(cfg.TranslationFile + ".tmp")
	if err != nil {
		return errors.Wrap(err, "creating tempfile")
	}

	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2) //nolint:gomnd

	if err = encoder.Encode(tf); err != nil {
		f.Close() //nolint:errcheck,gosec,revive // Short-lived fd-leak
		return errors.Wrap(err, "encoding translation file")
	}

	f.Close() //nolint:errcheck,gosec,revive // Short-lived fd-leak
	return errors.Wrap(os.Rename(cfg.TranslationFile+".tmp", cfg.TranslationFile), "moving file in place")
}
