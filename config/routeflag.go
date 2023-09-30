package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type routeFlag struct {
	routes []Route
}

func (f *routeFlag) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *routeFlag) Set(value string) error {
	route, err := parseRoute(value)
	if err != nil {
		return err
	}
	f.routes = append(f.routes, route)
	return nil
}

func parseRoute(value string) (Route, error) {
	source, target, found := strings.Cut(value, "=")
	if !found {
		return Route{}, errors.New("missing '='")
	}
	targetURL, err := url.Parse(target)
	if err != nil {
		return Route{}, err
	}
	if err := validateTarget(targetURL); err != nil {
		return Route{}, err
	}
	return NewRoute(source, targetURL), nil
}

func validateTarget(target *url.URL) error {
	switch target.Scheme {
	case "", "file":
		info, err := os.Stat(target.Path)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("%s: not a directory", target.Path)
		}
	}
	return nil
}
