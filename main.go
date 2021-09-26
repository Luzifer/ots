package main

import (
	"embed"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	http_helpers "github.com/Luzifer/go_helpers/v2/http"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg struct {
		Listen         string `flag:"listen" default:":3000" description:"IP/Port to listen on"`
		LogLevel       string `flag:"log-level" default:"info" description:"Set log level (debug, info, warning, error)"`
		StorageType    string `flag:"storage-type" default:"mem" description:"Storage to use for putting secrets to" validate:"nonzero"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Print version information and exit"`
	}

	product = "ots"
	version = "dev"
)

//go:embed frontend/*
var assets embed.FS

func init() {
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Error parsing CLI arguments: %s", err)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err == nil {
		log.SetLevel(l)
	} else {
		log.Fatalf("Invalid log level: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("%s %s\n", product, version)
		os.Exit(0)
	}
}

func main() {
	store, err := getStorageByType(cfg.StorageType)
	if err != nil {
		log.Fatalf("Could not initialize storage: %s", err)
	}
	api := newAPI(store)

	r := mux.NewRouter()
	r.Use(http_helpers.GzipHandler)

	api.Register(r.PathPrefix("/api").Subrouter())

	r.HandleFunc("/", handleIndex)
	r.PathPrefix("/").HandlerFunc(assetDelivery)

	log.Fatalf("HTTP server quit: %s", http.ListenAndServe(cfg.Listen, http_helpers.NewHTTPLogHandler(r)))
}

func assetDelivery(w http.ResponseWriter, r *http.Request) {
	assetName := r.URL.Path

	dot := strings.LastIndex(assetName, ".")
	if dot < 0 {
		// There are no assets with no dot in it
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	ext := assetName[dot:]
	assetData, err := assets.ReadFile(path.Join("frontend", assetName))
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(ext))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(assetData)
}

var (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
	cspHeader = strings.Join([]string{
		"default-src 'none'",
		"connect-src 'self'",
		"font-src 'self'",
		"img-src 'self'",
		"script-src 'self' 'unsafe-inline'",
		"style-src 'self' 'unsafe-inline'",
	}, ";")

	indexTpl *template.Template
)

func init() {
	source, err := assets.ReadFile("frontend/index.html")
	if err != nil {
		log.WithError(err).Fatal("frontend folder should contain index.html Go template")
	}
	indexTpl = template.Must(template.New("index.html").Funcs(tplFuncs).Parse(string(source)))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Xss-Protection", "1; mode=block")
	w.Header().Set("Content-Security-Policy", cspHeader)
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if err := indexTpl.Execute(w, struct {
		Version string
	}{
		Version: version,
	}); err != nil {
		http.Error(w, errors.Wrap(err, "executing template").Error(), http.StatusInternalServerError)
		return
	}
}
