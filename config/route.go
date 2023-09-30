package config

import (
	"fmt"
	"net/url"
	"strings"
)

type Route struct {
	Source string
	Target *url.URL
}

func NewRoute(source string, target *url.URL) Route {
	return Route{Source: strings.TrimRight(source, "/"), Target: target}
}

func (r Route) String() string {
	return fmt.Sprintf("%s=%s", r.Source, r.Target.String())
}
