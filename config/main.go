package config

import (
	"flag"
	"net"
)

func ReadConfig(args []string) (*Config, error) {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	configFile := flags.String("config", "", "path to configuration file")

	var addr *string
	flags.Func("listen", "address to listen on (default :8000)", func(value string) error {
		addr = &value
		_, err := net.ResolveTCPAddr("tcp", value)
		return err
	})

	var level *LogLevel = nil
	flags.Func("loglevel", "log level (info|warn|error)", func(value string) error {
		level = &LogLevel{}
		return level.UnmarshalText([]byte(value))
	})

	var routes Routes
	flags.Var(&routes, "route", `route specification (default "/=.")

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

	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}

	conf, err := ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	if addr != nil {
		conf.Addr = *addr
	}
	if level != nil {
		conf.LogLevel = *level
	}
	if len(routes.Proxy) > 0 || len(routes.Static) > 0 {
		conf.Routes.Proxy = routes.Proxy
		conf.Routes.Static = routes.Static
	}

	return conf, nil
}
