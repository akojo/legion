package file

import (
	"crypto/tls"
	"log/slog"
	"os"

	"github.com/akojo/legion/handler"
	"gopkg.in/yaml.v3"
)

type YAML struct {
	Listen string     `yaml:"listen"`
	Level  slog.Level `yaml:"loglevel"`
	Routes Routes     `yaml:"routes"`
	TLS    TLS        `yaml:"tls"`
}

type Config struct {
	Addr   string
	Level  slog.Level
	Routes []handler.Route
	TLS    *tls.Config
}

func ReadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	decoded := &YAML{}
	err = yaml.NewDecoder(f).Decode(&decoded)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Addr:  decoded.Listen,
		Level: decoded.Level,
	}

	routes, err := decoded.Routes.Get()
	if err != nil {
		return nil, err
	}
	if len(routes) > 0 {
		config.Routes = routes
	}

	config.TLS, err = decoded.TLS.GetConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}
