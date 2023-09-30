package config

import (
	"flag"
	"net"
	"net/url"
)

func ParseFlags() (*Config, error) {
	addr := ":8000"
	flag.Func("listen", "address to listen on (default :8000)", func(value string) error {
		addr = value
		_, err := net.ResolveTCPAddr("tcp", value)
		return err
	})

	quiet := flag.Bool("quiet", false, "disable request logging")

	var route routeFlag
	flag.Var(&route, "route", `route specification (default "/=.")

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
	- /api/pets/1 -> https://www.example.com/v1/pets/1`)

	flag.Parse()

	routes := route.routes
	if len(routes) == 0 {
		routes = []Route{NewRoute("/", &url.URL{Path: "."})}
	}

	return &Config{
		Addr:      addr,
		EnableLog: !*quiet,
		Routes:    routes,
	}, nil
}
