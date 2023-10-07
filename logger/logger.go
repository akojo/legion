package logger

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func Middleware(log *slog.Logger, next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(writer, r)

		log.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			strconv.Itoa(writer.status)+" "+r.Method+" "+r.URL.Path,
			slog.String("method", r.Method),
			slog.String("proto", r.Proto),
			slog.String("path", r.URL.Path),
			slog.String("address", r.Host),
			slog.Int("status", writer.status),
			slog.Duration("duration", time.Since(start)),
			slog.String("user_agent", r.Header.Get("User-Agent")))
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
