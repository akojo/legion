package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	handler http.Handler
	tls     *tls.Config
}

func New(handler http.Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) AddTLSCertificate(certfile, keyfile string) error {
	cert, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return err
	}
	if s.tls == nil {
		s.tls = &tls.Config{
			Certificates: []tls.Certificate{},
			NextProtos:   []string{"h2"},
		}
	}
	s.tls.Certificates = append(s.tls.Certificates, cert)
	return nil
}

func (s *Server) ListenAndServe(addr string) error {
	listener, err := s.listen(addr)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Handler:   s.handler,
		TLSConfig: s.tls,
	}

	quit := make(chan os.Signal, 1)
	shutdown := make(chan error)

	go func() {
		err := srv.Serve(listener)
		if !errors.Is(err, http.ErrServerClosed) {
			shutdown <- err
		}
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	case err := <-shutdown:
		return err
	}
}

func (s *Server) listen(addr string) (net.Listener, error) {
	if s.tls != nil {
		return tls.Listen("tcp", addr, s.tls)
	}
	return net.Listen("tcp", addr)
}
