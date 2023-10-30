package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func (h *Handler) WebServer(source, basedir string) error {
	basedir, err := ensureDir(basedir)
	if err != nil {
		return err
	}
	return h.addHandler(source, webServer(basedir))
}

func (h *Handler) ReverseProxy(source, URL string) error {
	target, err := url.Parse(URL)
	if err != nil {
		return err
	}
	return h.addHandler(source, reverseProxy(target))
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
	return strings.TrimSuffix(filepath.ToSlash(dirname), "/"), nil
}
