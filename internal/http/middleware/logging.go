package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type LoggingResponseWriter struct {
	ResponseWriter http.ResponseWriter
	StatusCode     int
}

func (lw *LoggingResponseWriter) Header() http.Header {
	return lw.ResponseWriter.Header()
}

func (lw *LoggingResponseWriter) WriteHeader(status int) {
	lw.StatusCode = status
	lw.ResponseWriter.WriteHeader(status)
}

func (lw *LoggingResponseWriter) Write(bytes []byte) (int, error) {
	if lw.StatusCode == 0 {
		lw.StatusCode = http.StatusOK
	}
	return lw.ResponseWriter.Write(bytes)
}

func NewLoggingMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lw := &LoggingResponseWriter{
				ResponseWriter: w,
				StatusCode:     http.StatusOK,
			}

			next.ServeHTTP(lw, r)

			duration := time.Since(start)

			fmt.Printf(
				"[%s] %s %s %d %s\n",
				duration,
				r.Method,
				r.URL.Path,
				lw.StatusCode,
				r.RemoteAddr,
			)
		})
	}
}
