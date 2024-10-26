package main

import (
	"crypto/rand"
	"embed"
	"encoding/base64"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	file_helpers "github.com/Luzifer/go_helpers/v2/file"
	http_helpers "github.com/Luzifer/go_helpers/v2/http"
	"github.com/Luzifer/ots/pkg/customization"
	"github.com/Luzifer/ots/pkg/metrics"
	"github.com/Luzifer/rconfig/v2"
)

const scriptNonceSize = 32

var (
	cfg struct {
		Customize      string `flag:"customize" default:"" description:"Customize-File to load"`
		Listen         string `flag:"listen" default:":3000" description:"IP/Port to listen on"`
		LogRequests    bool   `flag:"log-requests" default:"true" description:"Enable request logging"`
		LogLevel       string `flag:"log-level" default:"info" description:"Set log level (debug, info, warning, error)"`
		SecretExpiry   int64  `flag:"secret-expiry" default:"0" description:"Maximum expiry of the stored secrets in seconds"`
		StorageType    string `flag:"storage-type" default:"mem" description:"Storage to use for putting secrets to" validate:"nonzero"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Print version information and exit"`
		EnableTLS      bool   `flag:"enable-tls" default:"false" description:"Enable HTTPS/TLS"`
		CertFile       string `flag:"cert-file" default:"" description:"Path to the TLS certificate file"`
		KeyFile        string `flag:"key-file" default:"" description:"Path to the TLS private key file"`
	}

	assets   file_helpers.FSStack
	cust     customization.Customize
	indexTpl *template.Template

	version = "dev"
)

//go:embed frontend/*
var embeddedAssets embed.FS

func defaultCSP() http_helpers.CSP {
	c := http_helpers.CSP{}

	c.Add("base-uri", http_helpers.CSPSrcSelf)
	c.Add("default-src", http_helpers.CSPSrcNone)
	c.Add("connect-src", http_helpers.CSPSrcSelf)
	c.Add("font-src", http_helpers.CSPSrcSelf)
	c.Add("img-src", http_helpers.CSPSrcSelf)
	c.Add("img-src", http_helpers.CSPSrcSchemeData)
	c.Add("script-src", http_helpers.CSPSrcSelf)
	c.Add("style-src", http_helpers.CSPSrcSelf)

	return c
}

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

	if cust, err = customization.Load(cfg.Customize); err != nil {
		return errors.Wrap(err, "loading customizations")
	}

	frontendFS, err := fs.Sub(embeddedAssets, "frontend")
	if err != nil {
		return errors.Wrap(err, "creating sub-fs for assets")
	}
	assets = append(assets, frontendFS)

	if cust.OverlayFSPath != "" {
		assets = append(file_helpers.FSStack{os.DirFS(cust.OverlayFSPath)}, assets...)
	}

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	if cfg.VersionAndExit {
		logrus.WithField("version", version).Info("ots")
		os.Exit(0)
	}

	// Initialize metrics collector
	collector := metrics.New()

	// Initialize index template in order not to parse it multiple times
	source, err := assets.ReadFile("index.html")
	if err != nil {
		logrus.WithError(err).Fatal("frontend folder should contain index.html Go template")
	}
	indexTpl = template.Must(template.New("index.html").Funcs(tplFuncs).Parse(string(source)))

	// Initialize storage
	store, err := getStorageByType(cfg.StorageType)
	if err != nil {
		logrus.WithError(err).Fatal("initializing storage")
	}
	api := newAPI(store, collector)

	// Initialize server
	r := mux.NewRouter()

	api.Register(r.PathPrefix("/api").Subrouter())

	r.Handle("/metrics", metrics.Handler()).
		Methods(http.MethodGet).
		MatcherFunc(func(r *http.Request, _ *mux.RouteMatch) bool {
			return requestInSubnetList(r, cust.MetricsAllowedSubnets)
		})

	r.HandleFunc("/", handleIndex).
		Methods(http.MethodGet)
	r.PathPrefix("/").HandlerFunc(assetDelivery).
		Methods(http.MethodGet)

	var hdl http.Handler = r
	hdl = http_helpers.GzipHandler(hdl)
	if cfg.LogRequests {
		hdl = http_helpers.NewHTTPLogHandlerWithLogger(hdl, logrus.StandardLogger())
	}

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           hdl,
		ReadHeaderTimeout: time.Second,
	}

	// Start periodic stored metrics update (required for multi-instance
	// OTS hosting as other instances will create / delete secrets and
	// we need to keep up with that)
	go func() {
		for t := time.NewTicker(time.Minute); ; <-t.C {
			updateStoredSecretsCount(store, collector)
		}
	}()

	// Start server
	logrus.WithFields(logrus.Fields{
		"secret_expiry": time.Duration(cfg.SecretExpiry) * time.Second,
		"version":       version,
	}).Info("ots started")

	if cfg.EnableTLS {
		if cfg.CertFile == "" || cfg.KeyFile == "" {
			logrus.Fatal("TLS is enabled but cert-file or key-file is not provided")
		}
		logrus.Infof("Starting HTTPS server on %s", cfg.Listen)
		if err := server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile); err != nil {
			logrus.WithError(err).Fatal("HTTPS server quit unexpectedly")
		}
	} else {
		logrus.Infof("Starting HTTP server on %s", cfg.Listen)
		if err := server.ListenAndServe(); err != nil {
			logrus.WithError(err).Fatal("HTTP server quit unexpectedly")
		}
	}
}

func assetDelivery(w http.ResponseWriter, r *http.Request) {
	assetName := strings.TrimLeft(r.URL.Path, "/")

	dot := strings.LastIndex(assetName, ".")
	if dot < 0 {
		// There are no assets with no dot in it
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	ext := assetName[dot:]
	assetData, err := assets.ReadFile(assetName)
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(ext))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if _, err = w.Write(assetData); err != nil {
		logrus.WithError(err).Error("writing asset data")
	}
}

func handleIndex(w http.ResponseWriter, _ *http.Request) {
	inlineContentNonce := make([]byte, scriptNonceSize)
	if _, err := rand.Read(inlineContentNonce); err != nil {
		logrus.WithError(err).Error("generating script nonce")
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	inlineContentNonceStr := base64.StdEncoding.EncodeToString(inlineContentNonce)

	policy := defaultCSP()
	policy.Add("script-src", http_helpers.CSPSrcNonce(inlineContentNonceStr))
	policy.Add("style-src", http_helpers.CSPSrcNonce(inlineContentNonceStr))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Xss-Protection", "1; mode=block")
	w.Header().Set("Content-Security-Policy", policy.ToHeaderValue())
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if err := indexTpl.Execute(w, struct {
		Customize          customization.Customize
		InlineContentNonce string
		MaxSecretExpiry    int64
		Version            string
	}{
		Customize:          cust,
		InlineContentNonce: inlineContentNonceStr,
		MaxSecretExpiry:    cfg.SecretExpiry,
		Version:            version,
	}); err != nil {
		http.Error(w, errors.Wrap(err, "executing template").Error(), http.StatusInternalServerError)
		return
	}
}
