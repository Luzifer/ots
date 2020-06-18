package main

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go -modtime 1 -md5checksum ./frontend/...

import (
	"encoding/json"
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
		DisableCreateInterface bool   `flag:"disable-create-interface" default:"false" description:"Removes the interface for secret creation"`
		Listen                 string `flag:"listen" default:":3000" description:"IP/Port to listen on"`
		LogLevel               string `flag:"log-level" default:"info" description:"Set log level (debug, info, warning, error)"`
		StorageType            string `flag:"storage-type" default:"mem" description:"Storage to use for putting secrets to" validate:"nonzero"`
		VersionAndExit         bool   `flag:"version" default:"false" description:"Print version information and exit"`
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

	if strings.LastIndex(assetName, ".") < 0 {
		// There are no assets with no dot in it
		http.Error(res, "404 not found", http.StatusNotFound)
		return
	}

	ext := assetName[strings.LastIndex(assetName, "."):]
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

	vars := map[string]interface{}{
		"disableCreateInterface": cfg.DisableCreateInterface,
		"version":                version,
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

	opts, err := toJSONString(vars)
	if err != nil {
		log.WithError(err).Error("Unable to encode JSON var")
		http.Error(w, "Unable to encode options", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "const otsOptions = %s", opts)
}

func normalizeLang(lang string) string {
	return strings.ToLower(strings.Split(lang, "-")[0])
}

func toJSONString(in interface{}) (string, error) {
	b, err := json.Marshal(in)
	return string(b), err
}
