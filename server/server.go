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
	"syscall"

	"github.com/akojo/legion/config"
	"github.com/akojo/legion/logger"
)

type Server struct {
	config *config.Config
	mux    *http.ServeMux
}

func New(config *config.Config) (*Server, error) {
	srv := &Server{config: config, mux: http.NewServeMux()}
	for _, route := range config.Routes {
		if err := srv.AddRoute(route); err != nil {
			return nil, err
		}
	}
	return srv, nil
}

func (s *Server) AddRoute(route config.Route) error {
	handler, err := makeProxy(route.Target)
	if err != nil {
		return err
	}
	s.mux.Handle(route.Pattern(), http.StripPrefix(route.Prefix(), handler))
	return nil
}

func (s *Server) Run() error {
	handler := http.Handler(s.mux)
	if s.config.EnableLog {
		handler = logger.Middleware(handler)
	}

	server := http.Server{Addr: s.config.Addr, Handler: handler}
	quit := make(chan os.Signal, 1)
	hangup := make(chan error)

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("server closed")
		} else if err != nil {
			hangup <- err
		}
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		return server.Shutdown(context.Background())
	case err := <-hangup:
		return err
	}

}

func makeProxy(target *url.URL) (http.Handler, error) {
	switch target.Scheme {
	case "", "file":
		return makeFileHandler(target)
	case "http", "https":
		return makeHTTPHandler(target)
	default:
		return nil, fmt.Errorf("invalid scheme: %s", target.Scheme)
	}
}

func makeFileHandler(target *url.URL) (http.Handler, error) {
	proxy := &httputil.ReverseProxy{
		Rewrite:   func(pr *httputil.ProxyRequest) {},
		Transport: http.NewFileTransport(http.Dir(target.Path)),
	}
	return proxy, nil
}

func makeHTTPHandler(target *url.URL) (http.Handler, error) {
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
