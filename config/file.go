package config

import (
	"log/slog"
	"net/url"
	"os"

	"github.com/akojo/legion/config/file"
	"github.com/akojo/legion/server"
)

func ReadFile(filename string) (*server.Config, error) {
	conf := defaultConfig()
	if filename == "" {
		return conf, nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	err = file.NewConfig(f, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func defaultConfig() *server.Config {
	defaultRoute, _ := server.NewRoute("/", &url.URL{Path: "."})
	return &server.Config{
		Addr:     ":8000",
		LogLevel: slog.LevelInfo,
		Routes:   []server.Route{defaultRoute},
	}
}
