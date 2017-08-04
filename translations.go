package main

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/nicksnyder/go-i18n/i18n"
	log "github.com/sirupsen/logrus"
)

func init() {
	for _, filename := range AssetNames() {
		if !strings.HasPrefix(filename, "frontend/locale") || !strings.HasSuffix(filename, ".all.json") {
			continue
		}

		translationData, _ := Asset(filename)
		if err := i18n.ParseTranslationFileBytes(filename, translationData); err != nil {
			log.Fatalf("Unable to load translation data %q: %s", filename, err)
		}
	}
}

func getTFuncMap(r *http.Request) template.FuncMap {
	cookie, _ := r.Cookie("lang")

	cookieLang := ""
	if cookie != nil {
		cookieLang = cookie.Value
	}
	qpLang := r.URL.Query().Get("hl")
	acceptLang := r.Header.Get("Accept-Language")
	defaultLang := "en-US" // known valid language

	T, _ := i18n.Tfunc(cookieLang, qpLang, acceptLang, defaultLang)
	return template.FuncMap{
		"T": T,
	}
}
