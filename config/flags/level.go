package flags

import (
	"log/slog"

	"gopkg.in/yaml.v3"
)

type LogLevel slog.Level

func (l *LogLevel) String() string {
	return l.Level().String()
}

func (l *LogLevel) Set(value string) error {
	return (*slog.Level)(l).UnmarshalText([]byte(value))
}

func (l *LogLevel) UnmarshalYAML(value *yaml.Node) error {
	return l.Set(value.Value)
}

func (l *LogLevel) Level() slog.Level {
	return (slog.Level)(*l)
}
