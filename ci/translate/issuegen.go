package main

import (
	_ "embed"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

//go:embed issue.tpl.md
var issueTemplate string

func generateIssue(tf translationFile) error {
	fm := template.FuncMap{
		"English":             tplEnglish(tf),
		"MissingTranslations": tplMissingTranslations(tf),
		"Ping":                tplPing(tf),
	}

	tpl, err := template.New("issue").Funcs(fm).Parse(issueTemplate)
	if err != nil {
		return errors.Wrap(err, "parsing issue template")
	}

	f, err := os.Create(cfg.IssueFile)
	if err != nil {
		return errors.Wrap(err, "opening issue file")
	}
	defer f.Close() //nolint:errcheck // Short-lived fd-leak

	return errors.Wrap(tpl.Execute(f, tf), "executing issue template")
}

func tplEnglish(tf translationFile) func(string) any {
	return func(key string) any {
		return tf.Reference.Translations[key]
	}
}

func tplMissingTranslations(tf translationFile) func(string) []string {
	return func(lang string) []string {
		missing, _, _ := tf.Translations[lang].Translations.GetErrorKeys(tf.Reference.Translations)
		sort.Strings(missing)
		return missing
	}
}

func tplPing(tf translationFile) func(string) string {
	return func(lang string) string {
		if len(tf.Translations[lang].Translators) == 0 {
			return "No translators to ping for this language."
		}

		var pings []string
		for _, t := range tf.Translations[lang].Translators {
			pings = append(pings, "@"+t)
		}

		return strings.Join([]string{"Ping", strings.Join(pings, ", ")}, " ")
	}
}
