package logger

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Middleware(next http.Handler) http.HandlerFunc {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(writer, r)

		logger.Info(
			strconv.Itoa(writer.status)+" "+r.Method+" "+r.URL.Path,
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("proto", r.Proto),
				slog.String("path", r.URL.Path)),
			slog.Group("resp",
				slog.Int("status_code", writer.status)),
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
