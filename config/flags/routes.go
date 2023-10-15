package flags

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/akojo/legion/server"
)

type Routes []server.Route

func (f *Routes) String() string {
	return fmt.Sprintf("%v", []server.Route(*f))
}

func (f *Routes) Set(value string) error {
	route, err := parseRoute(value)
	if err != nil {
		return err
	}
	*f = append(*f, route)
	return nil
}

func parseRoute(value string) (server.Route, error) {
	source, target, found := strings.Cut(value, "=")
	if !found {
		return server.Route{}, errors.New("missing '='")
	}
	targetURL, err := url.Parse(filepath.ToSlash(target))
	if err != nil {
		return server.Route{}, err
	}
	if err := validateTarget(targetURL); err != nil {
		return server.Route{}, err
	}
	return server.NewRoute(source, targetURL)
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
