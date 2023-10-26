package config_test

import (
	"log/slog"
	"testing"

	"github.com/akojo/legion/config"
)

func TestDefaultConfig(t *testing.T) {
	conf := newConf(t)
	if conf.Addr != ":8000" {
		t.Errorf("Addr: want :8000, got %s", conf.Addr)
	}
	if conf.LogLevel.Level != slog.LevelInfo {
		t.Errorf("LogLevel: want INFO, got %v", conf.LogLevel)
	}
	if len(conf.Routes.Static) != 1 {
		t.Errorf("Routes: want 1, got %v", conf.Routes.Static)
	}
	want := config.StaticRoute{"/", "."}
	route := conf.Routes.Static[0]
	if route != want {
		t.Errorf("default: want %v, got %v", want, route)
	}
}

func TestListenFlag(t *testing.T) {
	conf := newConf(t, "-listen", ":3000")
	if conf.Addr != ":3000" {
		t.Errorf("Addr: want :3000, got %s", conf.Addr)
	}
}

func TestLogLevelFlag(t *testing.T) {
	conf := newConf(t, "-loglevel", "warn")
	if conf.LogLevel.Level != slog.LevelWarn {
		t.Errorf("LogLevel: want WARN, got %v", conf.LogLevel)
	}
}

func TestRouteFlag(t *testing.T) {
	conf := newConf(t, "-route", "/=/www")
	if len(conf.Routes.Static) != 1 {
		t.Errorf("Routes: want 1, got %v", conf.Routes.Static)
	}
	want := config.StaticRoute{"/", "/www"}
	route := conf.Routes.Static[0]
	if route != want {
		t.Errorf("default: want %v, got %v", want, route)
	}
}

func TestRouteFlagWithURL(t *testing.T) {
	conf := newConf(t, "-route", "/=http://example.com")
	if got := len(conf.Routes.Proxy); got != 1 {
		t.Errorf("Routes: want 1, got %v", got)
	}
	want := config.ProxyRoute{"/", "http://example.com"}
	route := conf.Routes.Proxy[0]
	if route != want {
		t.Errorf("default: want %v, got %v", want, route)
	}
}

func TestConfigFile(t *testing.T) {
	conf := newConf(t, "-config", "testdata/config.yml")
	if conf.Addr != ":80" {
		t.Errorf("Addr: want :80, got %s", conf.Addr)
	}

	if conf.LogLevel.Level != slog.LevelError {
		t.Errorf("LogLevel: want ERROR, got %s", conf.LogLevel)
	}

	if got := len(conf.TLS.Certificates); got != 1 {
		t.Errorf("TLS certs: want 1, got %d", got)
	}
	cert := config.Certificate{CertFile: "domain.crt", KeyFile: "domain.key"}
	if got := conf.TLS.Certificates[0]; got != cert {
		t.Errorf("TLS cert: want %s, got %s", cert, got)
	}

	static := config.StaticRoute{"/", "."}
	if got := len(conf.Routes.Static); got != 1 {
		t.Errorf("static routes: want 1, got %d", got)
	}
	if got := conf.Routes.Static[0]; got != static {
		t.Errorf("static route: want %s, got %s", static, got)
	}

	proxies := []config.ProxyRoute{
		{"/http", "http://example.com/"},
		{"/https", "https://example.com/"},
	}
	if got := len(conf.Routes.Proxy); got != 2 {
		t.Errorf("proxy routes: want 2, got %d", got)
	}
	for i, route := range conf.Routes.Proxy {
		if route != proxies[i] {
			t.Errorf("proxy route %d: want %s, got %s", i, proxies[i], route)
		}
	}
}

func TestOverrideAddress(t *testing.T) {
	conf := newConf(t,
		"-config", "testdata/config.yml",
		"-listen", ":3000",
	)
	if conf.Addr != ":3000" {
		t.Errorf("Addr: want :3000, got %s", conf.Addr)
	}
}

func TestOverrideLogLevel(t *testing.T) {
	conf := newConf(t,
		"-config", "testdata/config.yml",
		"-loglevel", "info",
	)
	if conf.LogLevel.Level != slog.LevelInfo {
		t.Errorf("LogLevel: want INFO, got %s", conf.LogLevel)
	}
}

func TestOverrideRoutes(t *testing.T) {
	conf := newConf(t,
		"-config", "testdata/config.yml",
		"-route", "/=.",
	)
	if got := len(conf.Routes.Proxy); got > 0 {
		t.Errorf("want 0 proxy routes, got %d", got)
	}
	if got := len(conf.Routes.Static); got != 1 {
		t.Errorf("want 1 static route, got %d", got)
	}
	want := config.StaticRoute{"/", "."}
	if got := conf.Routes.Static[0]; got != want {
		t.Errorf("route: want %v, got %v", want, got)
	}
}

func newConf(t *testing.T, args ...string) *config.Config {
	conf, err := config.ReadConfig(args)
	if err != nil {
		t.Fatal(err)
	}
	return conf
}
