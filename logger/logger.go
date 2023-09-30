package logger

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Transport struct {
	Transport http.RoundTripper
	logger    *slog.Logger
}

func NewLogger(transport http.RoundTripper) http.RoundTripper {
	return &Transport{
		Transport: transport,
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.MessageKey {
					return slog.Attr{}
				}
				return a
			},
		})),
	}
}

func (l *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := l.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	go l.write(time.Since(start), req, resp.StatusCode)

	return resp, err
}

func (l *Transport) write(duration time.Duration, req *http.Request, status int) {
	l.logger.Info(
		"",
		slog.Group("req",
			slog.String("method", req.Method),
			slog.String("proto", req.Proto),
			slog.String("path", req.URL.EscapedPath())),
		slog.Group("resp",
			slog.Int("status_code", status)),
		slog.Duration("duration", duration),
		slog.String("user_agent", req.Header.Get("User-Agent")))
}
