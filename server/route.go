package server

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
)

type Route struct {
	host   string
	path   string
	Target *url.URL
}

func NewRoute(source string, target *url.URL) (Route, error) {
	pathStart := strings.IndexRune(source, '/')
	if pathStart < 0 {
		return Route{}, errors.New("source path must start with a '/'")
	}
	return Route{
		host:   source[0:pathStart],
		path:   strings.TrimRight(source[pathStart:], "/"),
		Target: cleanPath(target),
	}, nil
}

func (r Route) String() string {
	return fmt.Sprintf("%s=%s", r.Pattern(), r.Target.String())
}

func (r Route) Pattern() string {
	return path.Join(r.host, r.path) + "/"
}

func (r Route) Prefix() string {
	return r.path
}

func cleanPath(u *url.URL) *url.URL {
	u.Path = strings.TrimRight(u.EscapedPath(), "/")
	return u
}
