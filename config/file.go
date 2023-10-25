package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadFile(filename string) (*Config, error) {
	conf := defaultConfig()
	if filename == "" {
		return conf, nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fileConf := &Config{}
	err = yaml.NewDecoder(f).Decode(fileConf)
	if err != nil {
		return nil, err
	}

	if fileConf.Addr != "" {
		conf.Addr = fileConf.Addr
	}
	conf.LogLevel = fileConf.LogLevel

	if len(fileConf.Routes.Static) > 0 {
		conf.Routes.Static = fileConf.Routes.Static
	}
	if len(fileConf.Routes.Proxy) > 0 {
		conf.Routes.Proxy = fileConf.Routes.Proxy
	}

	conf.TLS = fileConf.TLS

	return conf, nil
}

func defaultConfig() *Config {
	return &Config{
		Addr:     ":8000",
		LogLevel: LogLevel{slog.LevelInfo},
		Routes: Routes{
			Static: []StaticRoute{{Source: "/", Target: "."}},
		},
	}
}
