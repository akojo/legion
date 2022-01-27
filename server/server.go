package server

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func (s *Server) Route(pattern string, target *url.URL) error {
	var handler http.Handler
	var err error
	if target.Scheme == "" {
		handler, err = makeStatic(target.Path)
	} else {
		handler, err = makeProxy(target)
	}
	if err != nil {
		return err
	}
	pattern = strings.TrimRight(pattern, "/")
	prefix := strings.TrimLeftFunc(pattern, func(r rune) bool { return r != '/' })
	s.mux.Handle(pattern+"/", http.StripPrefix(prefix, handler))
	return nil
}

func (s *Server) Run(addr string, logEnabled bool) {
	handler := http.Handler(s.mux)
	if logEnabled {
		handler = logger.Middleware(handler)
	}
	server := http.Server{Addr: addr, Handler: handler}

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s", err)
		} else if err != nil {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("shutdown: %s", err)
	}
}

func makeStatic(root string) (http.Handler, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s: not a directory", root)
	}
	return http.FileServer(http.Dir(root)), nil
}

func makeProxy(target *url.URL) (http.Handler, error) {
	if !strings.HasPrefix(target.Scheme, "http") {
		return nil, fmt.Errorf("invalid scheme: %s", target.Scheme)
	}
	return httputil.NewSingleHostReverseProxy(target), nil
}
