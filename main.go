package main

import (
	"log/slog"
	"os"

	"github.com/akojo/legion/config"
	"github.com/akojo/legion/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	conf, err := config.ParseFlags()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	srv, err := server.New(conf)
	if err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}

	if err := srv.Run(); err != nil {
		slog.Error("server closed unexpectedly", "error", err)
		os.Exit(1)
	}
}
