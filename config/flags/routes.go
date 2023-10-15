package flags

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/akojo/legion/handler"
)

type Routes []handler.Route

func (f *Routes) String() string {
	return fmt.Sprintf("%v", []handler.Route(*f))
}

func (f *Routes) Set(value string) error {
	route, err := parseRoute(value)
	if err != nil {
		return err
	}
	*f = append(*f, route)
	return nil
}

func parseRoute(value string) (handler.Route, error) {
	source, target, found := strings.Cut(value, "=")
	if !found {
		return handler.Route{}, errors.New("missing '='")
	}
	targetURL, err := url.Parse(filepath.ToSlash(target))
	if err != nil {
		return handler.Route{}, err
	}
	if err := validateTarget(targetURL); err != nil {
		return handler.Route{}, err
	}
	return handler.NewRoute(source, targetURL)
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
