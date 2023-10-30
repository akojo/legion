package config

import (
	"errors"
	"fmt"
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
		r.Proxy = append(r.Proxy, ProxyRoute{source, target})
	} else {
		r.Static = append(r.Static, StaticRoute{source, target})
	}
	return nil
}
