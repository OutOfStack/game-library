package middleware

import (
	"expvar"
	"net/http"
)

var m = struct {
	req       *expvar.Int
	err       *expvar.Int
	serverErr *expvar.Int
}{
	req:       expvar.NewInt("requests"),
	err:       expvar.NewInt("errors"),
	serverErr: expvar.NewInt("server_errors"),
}

type statusCodeResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader writes status code to header
func (w *statusCodeResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Metrics adds metrics
func Metrics(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		rw := &statusCodeResponseWriter{ResponseWriter: w}

		next.ServeHTTP(w, r)

		m.req.Add(1)

		if rw.StatusCode >= 400 {
			m.err.Add(1)
		}

		if rw.StatusCode >= 500 {
			m.serverErr.Add(1)
		}
	}

	return http.HandlerFunc(h)
}
