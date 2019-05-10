package main

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go -modtime 1 -md5checksum ./frontend/...

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	http_helpers "github.com/Luzifer/go_helpers/http"
	"github.com/Luzifer/rconfig"
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
	r.PathPrefix("/").HandlerFunc(http_helpers.GzipFunc(assetDelivery))

	log.Fatalf("HTTP server quit: %s", http.ListenAndServe(cfg.Listen, http_helpers.NewHTTPLogHandler(r)))
}

func assetDelivery(res http.ResponseWriter, r *http.Request) {
	assetName := r.URL.Path
	if assetName == "/" {
		assetName = "/index.html"
	}

	ext := assetName[strings.LastIndex(assetName, "."):]
	assetData, err := Asset(path.Join("frontend", assetName))
	if err != nil {
		http.Error(res, "404 not found", http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", mime.TypeByExtension(ext))

	if assetName != "/index.html" {
		// Do not use template engine on other files than index.html
		res.Write(assetData)
		return
	}

	tpl, err := template.New(assetName).Funcs(addTranslateFunc(tplFuncs, r)).Parse(string(assetData))

	if err != nil {
		log.Errorf("Template for asset %q has an error: %s", assetName, err)
		return
	}

	tpl.Execute(res, map[string]interface{}{
		"version": version,
	})
}
