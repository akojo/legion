package handler

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func reverseProxy(target *url.URL) http.Handler {
	target.Path = strings.TrimRight(target.EscapedPath(), "/")
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			setURL(r.Out.URL, target)
			setHeaders(r)
		},
		Transport: http.DefaultTransport,
	}
}

func setURL(u *url.URL, target *url.URL) {
	u.Scheme = target.Scheme
	u.Host = target.Host
	u.Path, u.RawPath = target.Path+u.Path, target.Path+u.EscapedPath()
}

func setHeaders(r *httputil.ProxyRequest) {
	clientIP, _, err := net.SplitHostPort(r.In.RemoteAddr)
	if err == nil {
		prior := r.In.Header["X-Forwarded-For"]
		if len(prior) > 0 {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		r.Out.Header.Set("X-Forwarded-For", clientIP)
	} else {
		r.Out.Header.Del("X-Forwarded-For")
	}

	host := r.In.Header.Get("X-Forwarded-Host")
	if len(host) > 0 {
		r.Out.Host = host
	} else {
		r.Out.Host = r.In.Host
	}

	proto := r.In.Header.Get("X-Forwarded-Proto")
	if len(proto) > 0 {
		r.Out.Header.Set("X-Forwarded-Proto", proto)
	} else if r.In.TLS == nil {
		r.Out.Header.Set("X-Forwarded-Proto", "http")
	} else {
		r.Out.Header.Set("X-Forwarded-Proto", "https")
	}
}
