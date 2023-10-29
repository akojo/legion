package main

import (
	"log/slog"
	"os"

	"github.com/akojo/legion/config"
	"github.com/akojo/legion/handler"
	"github.com/akojo/legion/logger"
	"github.com/akojo/legion/server"
)

func main() {
	var logLevel = new(slog.LevelVar)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	conf, err := config.ReadConfig(os.Args[1:])
	if err != nil {
		Fatal("invalid configuration", err)
	}

	logLevel.Set(conf.LogLevel.Level)

	h := handler.New()
	for _, route := range conf.Routes.Static {
		err := h.FileServer(route.Source, route.Target)
		if err != nil {
			Fatal("invalid route", err)
		}
	}
	for _, route := range conf.Routes.Proxy {
		err := h.ReverseProxy(route.Source, route.Target)
		if err != nil {
			Fatal("invalid route", err)
		}
	}

	srv := server.New(logger.Middleware(slog.Default(), h))

	for _, c := range conf.TLS.Certificates {
		err = srv.AddTLSCertificate(c.CertFile, c.KeyFile)
		if err != nil {
			Fatal("invalid TLS config", err)
		}
	}

	err = srv.ListenAndServe(conf.Addr)
	if err != nil {
		Fatal("server closed unexpectedly", err)
	}
}

func Fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}
