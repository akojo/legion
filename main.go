package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rootDir := flag.String("dir", cwd, "`directory` to serve files from")
	listenAddr := flag.String("listen", ":8000", "`address` to listen on")
	quiet := flag.Bool("quiet", false, "disable request logging")
	flag.Parse()

	handler := http.FileServer(http.Dir(*rootDir))
	if !*quiet {
		handler = logRequest(handler)
	}

	server := http.Server{Addr: *listenAddr, Handler: handler}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("shutdown: %s\n", err)
	}
}
