package main

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go -modtime 1 -md5checksum ./frontend/...

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
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
	api.Register(r.PathPrefix("/api").Subrouter())
	r.HandleFunc("/vars.js", handleVars)
	r.PathPrefix("/").HandlerFunc(http_helpers.GzipFunc(assetDelivery))

	log.Fatalf("HTTP server quit: %s", http.ListenAndServe(cfg.Listen, http_helpers.NewHTTPLogHandler(r)))
}

func assetDelivery(res http.ResponseWriter, r *http.Request) {
	assetName := r.URL.Path
	if assetName == "/" {
		assetName = "/index.html"
	}

	dot := strings.LastIndex(assetName, ".")
	if dot < 0 {
		// There are no assets with no dot in it
		http.Error(res, "404 not found", http.StatusNotFound)
		return
	}

	ext := assetName[dot:]
	assetData, err := Asset(path.Join("frontend", assetName))
	if err != nil {
		http.Error(res, "404 not found", http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", mime.TypeByExtension(ext))
	res.Write(assetData)
}

func handleVars(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("lang")

	cookieLang := ""
	if cookie != nil {
		cookieLang = cookie.Value
	}
	acceptLang := r.Header.Get("Accept-Language")
	defaultLang := "en" // known valid language

	vars := map[string]string{
		"version": version,
	}

	switch {
	case cookieLang != "":
		vars["locale"] = normalizeLang(cookieLang)
	case acceptLang != "":
		vars["locale"] = normalizeLang(strings.Split(acceptLang, ",")[0])
	default:
		vars["locale"] = defaultLang
	}

	w.Header().Set("Content-Type", "application/javascript")
	for k, v := range vars {
		fmt.Fprintf(w, "var %s = %q\n", k, v)
	}
}

func normalizeLang(lang string) string {
	return strings.ToLower(strings.Split(lang, "-")[0])
}
