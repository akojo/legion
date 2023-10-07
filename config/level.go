package config

import (
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"
)

var levelNames = map[slog.Level]string{
	slog.LevelInfo:  "debug",
	slog.LevelWarn:  "warn",
	slog.LevelError: "error",
}

type LogLevel slog.Level

func (l *LogLevel) String() string {
	return levelNames[slog.Level(*l)]
}

func (l *LogLevel) Set(value string) error {
	level, err := logLevel(value)
	if err != nil {
		return err
	}
	*l = LogLevel(level)
	return nil
}

func (l *LogLevel) UnmarshalYAML(value *yaml.Node) error {
	return l.Set(value.Value)
}

func logLevel(level string) (slog.Level, error) {
	switch level {
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid loglevel '%s'", level)
	}
}
