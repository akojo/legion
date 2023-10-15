package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func New(routes []Route) (http.Handler, error) {
	mux := http.NewServeMux()
	for _, route := range routes {
		handler, err := makeProxy(route.Target)
		if err != nil {
			return nil, err
		}
		mux.Handle(route.Pattern(), http.StripPrefix(route.Prefix(), handler))
	}
	return mux, nil
}

func makeProxy(target *url.URL) (http.Handler, error) {
	switch target.Scheme {
	case "", "file":
		return makeFileHandler(target)
	case "http", "https":
		return makeHTTPHandler(target)
	default:
		return nil, fmt.Errorf("invalid scheme: %s", target.Scheme)
	}
}

func makeFileHandler(target *url.URL) (http.Handler, error) {
	proxy := &httputil.ReverseProxy{
		Rewrite:   func(pr *httputil.ProxyRequest) {},
		Transport: http.NewFileTransport(http.Dir(target.Path)),
	}
	return proxy, nil
}

func makeHTTPHandler(target *url.URL) (http.Handler, error) {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.Out.Header["X-Forwarded-For"] = pr.In.Header["X-Forwarded-For"]
			pr.SetXForwarded()
			pr.SetURL(target)
			pr.Out.Host = pr.In.Host
		},
		Transport: http.DefaultTransport,
	}
	return proxy, nil
}
