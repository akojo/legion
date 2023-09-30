package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/akojo/legion/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	var routes routeFlags
	cwd, err := os.Getwd()
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		os.Exit(1)
	}
	listenAddr := flag.String("listen", ":8000", "`address` to listen on")
	quiet := flag.Bool("quiet", false, "disable request logging")
	flag.Var(&routes, "route", fmt.Sprintf(`route specification (default "/=%s")

Routes are specified with <source>=<target>. -route option can be
specified multiple times.

<source> can be either
    - path, e.g. /api
    - path prefixed by a hostname, e.g. www.example.com/api
In the latter case the path matches a request only when request Host:
header matches given hostname.

<target> can be either
    - local filesystem path, e.g. /var/www/html
    - URL to proxy requests to, e.g. www.example.com/api/v1
In either case source path is first stripped from incoming requests and
the result is appended to target.

As an example, given options
    -route /api=https://www.example.com/v1 -route /=/var/www/html
incoming paths map to actual requests:
	- /index.html -> /var/www/html/index.html
	- /api/pets/1 -> https://www.example.com/v1/pets/1`, cwd))
	flag.Parse()

	if len(routes) == 0 {
		routes = routeFlags{{path: "/", target: &url.URL{Path: cwd}}}
	}

	srv := server.New()
	for _, route := range routes {
		if err := srv.AddRoute(route.path, route.target, !*quiet); err != nil {
			slog.Error(fmt.Sprintf("%s", err))
			os.Exit(1)
		}
	}
	srv.Run(*listenAddr)
}
