package file

import (
	"net/url"

	"github.com/akojo/legion/handler"
)

type Routes struct {
	Static []Route `yaml:"static"`
	Proxy  []Route `yaml:"proxy"`
}

type Route struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

func (r Routes) Get() ([]handler.Route, error) {
	result := []handler.Route{}
	for _, r := range append(r.Static, r.Proxy...) {
		targetURL, err := url.Parse(r.Target)
		if err != nil {
			return nil, err
		}
		route, err := handler.NewRoute(r.Source, targetURL)
		if err != nil {
			return nil, err
		}
		result = append(result, route)
	}
	return result, nil
}
