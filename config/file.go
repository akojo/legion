package config

import (
	"io"
	"log/slog"
	"net/url"
	"os"

	"github.com/akojo/legion/server"
	"gopkg.in/yaml.v3"
)

type YAML struct {
	Listen string   `yaml:"listen"`
	Level  LogLevel `yaml:"loglevel"`
	Routes Routes   `yaml:"routes"`
}

type Routes struct {
	Static []Route `yaml:"static"`
	Proxy  []Route `yaml:"proxy"`
}

type Route struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

func ReadFile(filename string) (*server.Config, error) {
	conf := defaultConfig()
	if filename == "" {
		return conf, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	err = yamlConfig(file, conf)
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

func yamlConfig(r io.Reader, conf *server.Config) error {
	var c YAML
	err := yaml.NewDecoder(r).Decode(&c)
	if err != nil {
		return err
	}

	conf.Addr = c.Listen
	conf.LogLevel = slog.Level(c.Level)

	conf.Routes = []server.Route{}
	for _, r := range append(c.Routes.Static, c.Routes.Proxy...) {
		targetURL, err := url.Parse(r.Target)
		if err != nil {
			return err
		}
		route, err := server.NewRoute(r.Source, targetURL)
		if err != nil {
			return err
		}
		conf.Routes = append(conf.Routes, route)
	}

	return nil
}
