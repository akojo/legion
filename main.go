package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/akojo/legion/config"
	"github.com/akojo/legion/handler"
	"github.com/akojo/legion/logger"
)

func main() {
	var logLevel = new(slog.LevelVar)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	conf, err := config.ReadConfig(os.Args[1:])
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	logLevel.Set(conf.LogLevel.Level)

	h := handler.New()
	for _, route := range conf.Routes.Static {
		err := h.FileServer(route.Source, route.Target)
		if err != nil {
			slog.Error("invalid route", "error", err)
			os.Exit(1)
		}
	}
	for _, route := range conf.Routes.Proxy {
		err := h.ReverseProxy(route.Source, route.Target)
		if err != nil {
			slog.Error("invalid route", "error", err)
			os.Exit(1)
		}
	}

	tlsConfig, err := makeTLSConfig(conf.TLS)
	if err != nil {
		slog.Error("invalid TLS config", "error", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Handler:   logger.Middleware(slog.Default(), h),
		TLSConfig: tlsConfig,
	}

	if err = listenAndServe(srv, conf.Addr); err != nil {
		slog.Error("server closed unexpectedly", "error", err)
		os.Exit(1)
	}
}

func listenAndServe(srv *http.Server, addr string) error {
	listener, err := listen(addr, srv.TLSConfig)
	if err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	hangup := make(chan error)

	go func() {
		err := srv.Serve(listener)
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("server closed")
		} else if err != nil {
			hangup <- err
		}
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		return srv.Shutdown(context.Background())
	case err := <-hangup:
		return err
	}
}

func listen(addr string, tlsConfig *tls.Config) (net.Listener, error) {
	if tlsConfig != nil {
		return tls.Listen("tcp", addr, tlsConfig)
	}
	return net.Listen("tcp", addr)
}

func makeTLSConfig(t config.TLS) (*tls.Config, error) {
	if len(t.Certificates) == 0 {
		return nil, nil
	}

	certs := make([]tls.Certificate, 0)
	for _, c := range t.Certificates {
		cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	return &tls.Config{
		Certificates: certs,
		NextProtos:   []string{"h2"},
	}, nil
}
