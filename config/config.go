package config

import (
	"log/slog"
)

type Config struct {
	Addr     string   `yaml:"listen"`
	LogLevel LogLevel `yaml:"loglevel"`
	Routes   Routes   `yaml:"routes"`
	TLS      TLS      `yaml:"tls"`
}

type LogLevel struct {
	slog.Level
}

func (l *LogLevel) Set(value string) error {
	return l.UnmarshalText([]byte(value))
}

type Routes struct {
	Static []StaticRoute `yaml:"static"`
	Proxy  []ProxyRoute  `yaml:"proxy"`
}

type StaticRoute struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type ProxyRoute struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type TLS struct {
	Certificates []Certificate `yaml:"certificates"`
}

type Certificate struct {
	CertFile string `yaml:"certfile"`
	KeyFile  string `yaml:"keyfile"`
}
