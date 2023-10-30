package handler

import (
	"net/http"
	"os"
	"path"
	"strings"
)

type webServer string

func (root webServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const indexPage = "/index.html"

	if strings.HasSuffix(r.URL.Path, indexPage) {
		redirect(w, r, "./")
		return
	}

	fpath := string(root) + r.URL.Path

	f, err := os.Open(fpath)
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if !info.IsDir() {
		http.ServeContent(w, r, fpath, info.ModTime(), f)
		return
	}

	ipath := strings.TrimSuffix(fpath, "/") + indexPage
	ifile, err := os.OpenFile(ipath, os.O_RDONLY, 0)
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	defer ifile.Close()

	if r.URL.Path[len(r.URL.Path)-1] != '/' {
		redirect(w, r, path.Base(fpath)+"/")
		return
	}

	info, err = ifile.Stat()
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeContent(w, r, "", info.ModTime(), ifile)
}

func redirect(w http.ResponseWriter, r *http.Request, path string) {
	w.Header().Set("Location", path)
	w.WriteHeader(http.StatusMovedPermanently)
}
