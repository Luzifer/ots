package http

import (
	"log"
	"net/http"
	"time"

	"github.com/Luzifer/go_helpers/accessLogger"
)

type HTTPLogHandler struct {
	Handler http.Handler
}

func NewHTTPLogHandler(h http.Handler) http.Handler {
	return HTTPLogHandler{Handler: h}
}

func (l HTTPLogHandler) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ares := accessLogger.New(res)

	l.Handler.ServeHTTP(ares, r)

	log.Printf("%s - \"%s %s\" %d %d \"%s\" \"%s\" %s",
		r.RemoteAddr,
		r.Method,
		r.URL.Path,
		ares.StatusCode,
		ares.Size,
		r.Header.Get("Referer"),
		r.Header.Get("User-Agent"),
		time.Since(start),
	)
}
