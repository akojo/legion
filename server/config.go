package server

import (
	"crypto/tls"
	"log/slog"
)

type Config struct {
	Addr     string
	LogLevel slog.Level
	Routes   []Route
	TLS      *tls.Config
}
