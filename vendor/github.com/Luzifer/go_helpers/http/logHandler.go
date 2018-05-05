package http

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Luzifer/go_helpers/accessLogger"
)

type HTTPLogHandler struct {
	Handler          http.Handler
	TrustedIPHeaders []string
}

func NewHTTPLogHandler(h http.Handler) http.Handler {
	return HTTPLogHandler{
		Handler:          h,
		TrustedIPHeaders: []string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"},
	}
}

func (l HTTPLogHandler) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ares := accessLogger.New(res)

	l.Handler.ServeHTTP(ares, r)

	log.Printf("%s - \"%s %s\" %d %d \"%s\" \"%s\" %s",
		l.findIP(r),
		r.Method,
		r.URL.Path,
		ares.StatusCode,
		ares.Size,
		r.Header.Get("Referer"),
		r.Header.Get("User-Agent"),
		time.Since(start),
	)
}

func (l HTTPLogHandler) findIP(r *http.Request) string {
	remoteAddr := strings.SplitN(r.RemoteAddr, ":", 2)[0]

	for _, hdr := range l.TrustedIPHeaders {
		if value := r.Header.Get(hdr); value != "" {
			return strings.SplitN(value, ",", 2)[0]
		}
	}

	return remoteAddr
}
