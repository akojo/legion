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
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return route{}, errors.New("missing '='")
	}
	url, err := url.Parse(strings.Join(parts[1:], "="))
	if err != nil {
		return route{}, err
	}
	return route{path: parts[0], target: url}, nil
}
