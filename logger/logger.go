package logger

import (
	"io"
	"log"
	"net/http"
	"time"
)

type timestampWriter struct {
	wrapped io.Writer
}

func (w *timestampWriter) Write(p []byte) (int, error) {
	prefix := time.Now().UTC().Format(time.RFC3339Nano) + " "
	if _, err := w.wrapped.Write([]byte(prefix)); err != nil {
		return 0, err
	}
	return w.wrapped.Write(p)
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(&timestampWriter{log.Writer()})
}

func Middleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &writer{rw: w, code: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(resp, req)
		log.Printf("%d %s %s %s %d B %.3f ms %s",
			resp.code,
			req.Method,
			req.Proto,
			req.URL.EscapedPath(),
			resp.bytes,
			float64(time.Since(start).Microseconds())/1000.0,
			req.Header.Get("User-Agent"))
	}
}

type writer struct {
	rw    http.ResponseWriter
	bytes int64
	code  int
}

func (w *writer) Header() http.Header {
	return w.rw.Header()
}

func (w *writer) Write(buf []byte) (int, error) {
	w.bytes += int64(len(buf))
	return w.rw.Write(buf)
}

func (w *writer) WriteHeader(statusCode int) {
	w.code = statusCode
	w.rw.WriteHeader(statusCode)
}
