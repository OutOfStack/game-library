package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Logger logs incoming HTTP requests and their responses using zap.Logger
func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// wrap the ResponseWriter to capture the status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			// log the request and response details
			statusCode := ww.Status()
			logFields := []zap.Field{
				zap.Int("status", statusCode),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			}

			// error on 5xx
			if statusCode >= 500 {
				logger.Error("Server error", logFields...)
				return
			}

			// warn on 401 or 403
			if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
				logger.Info("Unauthorized access attempt", logFields...)
				return
			}

			// info on other
			logger.Info("HTTP Request",
				append(logFields,
					zap.Int("bytes", ww.BytesWritten()),
				)...)
		})
	}
}
