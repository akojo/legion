package config

import (
	"log/slog"
	"net/url"

	"github.com/akojo/legion/config/file"
	"github.com/akojo/legion/handler"
)

func ReadFile(filename string) (*Config, error) {
	conf := defaultConfig()
	if filename == "" {
		return conf, nil
	}

	fileConf, err := file.ReadConfig(filename)
	if err != nil {
		return nil, err
	}

	if fileConf.Addr != "" {
		conf.Addr = fileConf.Addr
	}
	conf.LogLevel = fileConf.Level

	if len(fileConf.Routes) > 0 {
		conf.Routes = fileConf.Routes
	}
	conf.TLS = fileConf.TLS

	return conf, nil
}

func defaultConfig() *Config {
	defaultRoute, _ := handler.NewRoute("/", &url.URL{Path: "."})
	return &Config{
		Addr:     ":8000",
		LogLevel: slog.LevelInfo,
		Routes:   []handler.Route{defaultRoute},
	}
}
