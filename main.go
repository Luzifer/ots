package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Luzifer/rconfig"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

	log.Fatalf("HTTP server quit: %s", http.ListenAndServe(cfg.Listen, r))
}
