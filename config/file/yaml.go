package file

import (
	"io"
	"log/slog"

	"github.com/akojo/legion/server"
	"gopkg.in/yaml.v3"
)

type YAML struct {
	Listen string     `yaml:"listen"`
	Level  slog.Level `yaml:"loglevel"`
	Routes Routes     `yaml:"routes"`
	TLS    TLS        `yaml:"tls"`
}

func NewConfig(r io.Reader, conf *server.Config) error {
	var c YAML
	err := yaml.NewDecoder(r).Decode(&c)
	if err != nil {
		return err
	}

	conf.Addr = c.Listen
	conf.LogLevel = slog.Level(c.Level)

	conf.Routes, err = c.Routes.Get()
	if err != nil {
		return err
	}
	conf.TLS, err = c.TLS.GetConfig()
	if err != nil {
		return err
	}
	return nil
}
