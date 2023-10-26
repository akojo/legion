package handler

import (
	"fmt"
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
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.Out.Header["X-Forwarded-For"] = pr.In.Header["X-Forwarded-For"]
			pr.SetXForwarded()
			pr.SetURL(target)
			pr.Out.Host = pr.In.Host
		},
		Transport: http.DefaultTransport,
	}
	return h.addHandler(source, handler)
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
