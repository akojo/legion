package main

import (
	"log/slog"
	"os"

	"github.com/akojo/legion/config"
	"github.com/akojo/legion/server"
)

func main() {
	var logLevel = new(slog.LevelVar)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	conf, err := config.ReadConfig()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	logLevel.Set(slog.Level(conf.LogLevel))

	srv, err := server.New(conf)
	if err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}

	err = srv.Run(slog.Default())
	if err != nil {
		slog.Error("server closed unexpectedly", "error", err)
		os.Exit(1)
	}
}
