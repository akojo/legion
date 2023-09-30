package config

import (
	"errors"
	"fmt"
	"net/url"
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
	return NewRoute(source, targetURL), nil
}
