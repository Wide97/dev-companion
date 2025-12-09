package middleware

import "net/http"

type LoggingResponseWriter struct {
	ResponseWriter http.ResponseWriter
	StatusCode     int
}

func (lw *LoggingResponseWriter) WriteHeader(status int) {
	lw.StatusCode = status
	lw.ResponseWriter.WriteHeader(status)

}

func (lw *LoggingResponseWriter) Write(bytes []byte) (int, error) {
	if lw.StatusCode == 0 {
		lw.StatusCode = 200
	}
	n, err := lw.ResponseWriter.Write(bytes)
	if err != nil {
		return n, err
	}

	return n, nil

}

func NewLoggingMiddleware() func(next http.Handler) http.Handler {

}
