package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type route struct {
	path   string
	target *url.URL
}

type routeFlags []route

func (f *routeFlags) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *routeFlags) Set(value string) error {
	route, err := parseRoute(value)
	if err != nil {
		return err
	}
	*f = append(*f, route)
	return nil
}

func parseRoute(value string) (route, error) {
	source, target, found := strings.Cut(value, "=")
	if !found {
		return route{}, errors.New("missing '='")
	}
	url, err := url.Parse(target)
	if err != nil {
		return route{}, err
	}
	return route{path: source, target: url}, nil
}
