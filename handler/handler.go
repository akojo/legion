package handler

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Handler struct {
	*http.ServeMux
}

func New() *Handler {
	return &Handler{ServeMux: http.NewServeMux()}
}

func (h *Handler) FileServer(source, dirname string) error {
	dirname, err := ensureDir(dirname)
	if err != nil {
		return err
	}
	return h.addHandler(source, http.FileServer(http.Dir(dirname)))
}

func (h *Handler) ReverseProxy(source, URL string) error {
	target, err := url.Parse(URL)
	if err != nil {
		return err
	}
	target.Path = strings.TrimRight(target.EscapedPath(), "/")

	handler := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			setURL(r.Out.URL, target)
			setHeaders(r)
		},
		Transport: http.DefaultTransport,
	}
	return h.addHandler(source, handler)
}

func setURL(u *url.URL, target *url.URL) {
	u.Scheme = target.Scheme
	u.Host = target.Host
	u.Path, u.RawPath = u.Path+target.Path, u.EscapedPath()+target.Path
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

func (h *Handler) addHandler(source string, handler http.Handler) error {
	pathStart := strings.Index(source, "/")
	if pathStart < 0 {
		return fmt.Errorf("%s: source path must start with '/'", source)
	}
	pattern := strings.TrimRight(source, "/") + "/"
	prefix := strings.TrimRight(source[pathStart:], "/")
	h.Handle(pattern, http.StripPrefix(prefix, handler))
	return nil
}

func ensureDir(dirname string) (string, error) {
	info, err := os.Stat(dirname)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s: not a directory", dirname)
	}
	return dirname, nil
}
