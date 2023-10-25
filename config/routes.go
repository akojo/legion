package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func (r *Routes) String() string {
	return fmt.Sprintf("%#v", r)
}

func (r *Routes) Set(value string) error {
	source, target, found := strings.Cut(value, "=")
	if !found {
		return errors.New("missing '='")
	}
	if strings.HasPrefix(target, "http:") || strings.HasPrefix(target, "https:") {
		return r.proxyRoute(source, target)
	}
	return r.staticRoute(source, target)
}

func (r *Routes) proxyRoute(source, target string) error {
	if _, err := url.Parse(target); err != nil {
		return errors.New("invalid URL")
	}
	r.Proxy = append(r.Proxy, ProxyRoute{source, target})
	return nil
}

func (r *Routes) staticRoute(source, target string) error {
	if strings.HasPrefix(target, "file:") {
		if _, err := url.Parse(target); err != nil {
			return errors.New("invalid file URL")
		}
	}
	r.Static = append(r.Static, StaticRoute{source, target})
	return nil
}
