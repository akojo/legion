package server

import "log/slog"

type Config struct {
	Addr     string
	LogLevel slog.Level
	Routes   []Route
}
