package accessLogger

import (
	"fmt"
	"net/http"
	"strconv"
)

type AccessLogResponseWriter struct {
	StatusCode int
	Size       int

	http.ResponseWriter
}

func New(res http.ResponseWriter) *AccessLogResponseWriter {
	return &AccessLogResponseWriter{
		StatusCode:     200,
		Size:           0,
		ResponseWriter: res,
	}
}

func (a *AccessLogResponseWriter) Write(out []byte) (int, error) {
	s, err := a.ResponseWriter.Write(out)
	a.Size += s
	return s, err
}

func (a *AccessLogResponseWriter) WriteHeader(code int) {
	a.StatusCode = code
	a.ResponseWriter.WriteHeader(code)
}

func (a *AccessLogResponseWriter) HTTPResponseType() string {
	return fmt.Sprintf("%cxx", strconv.FormatInt(int64(a.StatusCode), 10)[0])
}
