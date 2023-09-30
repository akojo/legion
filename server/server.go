package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/akojo/legion/logger"
)

type Server struct {
	mux *http.ServeMux
}

func New() *Server {
	return &Server{mux: http.NewServeMux()}
}

func (s *Server) AddRoute(pattern string, target *url.URL, log bool) error {
	proxy, err := makeProxy(target)
	if err != nil {
		return err
	}
	if log {
		proxy.Transport = logger.NewLogger(proxy.Transport)
	}
	pattern = strings.TrimRight(pattern, "/")
	prefix := strings.TrimLeftFunc(pattern, func(r rune) bool { return r != '/' })
	s.mux.Handle(pattern+"/", http.StripPrefix(prefix, proxy))
	return nil
}

func (s *Server) Run(addr string) {
	server := http.Server{Addr: addr, Handler: http.Handler(s.mux)}
	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info(fmt.Sprintf("listen: %s", err))
		} else if err != nil {
			slog.Error(fmt.Sprintf("listen: %s", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		slog.Warn(fmt.Sprintf("shutdown: %s", err))
	}
}

func makeProxy(target *url.URL) (*httputil.ReverseProxy, error) {
	switch target.Scheme {
	case "", "file":
		return makeFileHandler(target)
	case "http", "https":
		return makeHTTPHandler(target)
	default:
		return nil, fmt.Errorf("invalid scheme: %s", target.Scheme)
	}
}

func makeFileHandler(target *url.URL) (*httputil.ReverseProxy, error) {
	root := target.Path
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s: not a directory", root)
	}
	proxy := &httputil.ReverseProxy{
		Rewrite:   func(pr *httputil.ProxyRequest) {},
		Transport: http.NewFileTransport(http.Dir(root)),
	}
	return proxy, nil
}

func makeHTTPHandler(target *url.URL) (*httputil.ReverseProxy, error) {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.Out.Header["X-Forwarded-For"] = pr.In.Header["X-Forwarded-For"]
			pr.SetXForwarded()
			pr.SetURL(target)
			pr.Out.Host = pr.In.Host
		},
		Transport: http.DefaultTransport,
	}
	return proxy, nil
}
