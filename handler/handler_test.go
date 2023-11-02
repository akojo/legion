package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akojo/legion/handler"
)

var emptyHeader = http.Header{}

type BenchmarkResponseWriter struct{}

func (w *BenchmarkResponseWriter) Header() http.Header {
	return emptyHeader
}

func (w *BenchmarkResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (w *BenchmarkResponseWriter) WriteHeader(statusCode int) {
}

func NewBenchmarkResponseWriter() *BenchmarkResponseWriter {
	return &BenchmarkResponseWriter{}
}

func TestRedirects(t *testing.T) {
	type test struct {
		path string
		want string
	}
	tests := []test{
		{path: "/index.html", want: "./"},
		{path: "/.", want: "/"},
		{path: "/..", want: "/"},
		{path: "//", want: "/"},
		{path: "/subdir/index.html", want: "./"},
		{path: "/subdir", want: "subdir/"},
		{path: "/./subdir", want: "/subdir"},
		{path: "/../../subdir", want: "/subdir"},
		{path: "/foo/../subdir", want: "/subdir"},
	}

	h := makeFileserver(t, "/", "testdata/html")

	for _, tc := range tests {
		resp := GET(h, tc.path)

		if status := resp.Result().StatusCode; status != 301 {
			t.Errorf("%s: status: want 301, got %d", tc.path, status)
		}
		if location := resp.Result().Header.Get("Location"); location != tc.want {
			t.Errorf("%s: location: want %#v, got %#v", tc.path, tc.want, location)
		}
	}
}

func TestGetPages(t *testing.T) {
	type test struct {
		path  string
		title string
	}
	tests := []test{
		{path: "/", title: "Main Page"},
		{path: "/subdir/", title: "Subdirectory"},
		{path: "/subpage.html", title: "Subpage"},
	}

	for _, tc := range tests {
		resp := GET(makeFileserver(t, "/", "testdata/html"), tc.path)

		if status := resp.Result().StatusCode; status != 200 {
			t.Fatalf("%s: status: want 200, got %d", tc.path, status)
		}
		if title := readTitle(t, resp.Result()); title != tc.title {
			t.Errorf("%s: want %#v, got %#v", tc.path, tc.title, title)
		}
	}
}

func TestNonExistentTargetDir(t *testing.T) {
	h := handler.New()
	err := h.FileServer("/", "./nosuchdir")
	if err == nil {
		t.Error("expect error")
	} else if !strings.Contains(err.Error(), "nosuchdir") {
		t.Errorf("expect %#v to contain 'nosuchdir'", err.Error())
	}
}

func TestTargetNotADir(t *testing.T) {
	h := handler.New()
	err := h.FileServer("/", "testdata/html/index.html")
	if err == nil {
		t.Error("expect error")
	} else if !strings.Contains(err.Error(), "index.html") {
		t.Errorf("expect %#v to contain 'index.html'", err.Error())
	}
}

func TestInvalidSource(t *testing.T) {
	h := handler.New()
	err := h.FileServer("invalidsource", ".")
	if err == nil {
		t.Error("expect error")
	} else if !strings.Contains(err.Error(), "invalidsource") {
		t.Errorf("expect %#v to contain 'invalidsource'", err.Error())
	}
}

func TestProxy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello from proxy")
	}))
	defer server.Close()

	resp := GET(makeReverseProxy(t, "/", server.URL), "/")
	if status := resp.Result().StatusCode; status != 200 {
		t.Errorf("want 200, got %d", status)
	}
	if body := readBody(t, resp.Result()); body != "hello from proxy" {
		t.Errorf("want 'hello from proxy', got %#v", body)
	}
}

func TestProxyHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := "host.example.com"
		if r.Host != host {
			t.Errorf("Host: want %#v, got %#v", host, r.Host)
		}
		if got := r.Header.Get("X-Forwarded-Proto"); got != "http" {
			t.Errorf("X-Forwarded-Proto: want 'http', got %#v", got)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	resp := GET(makeReverseProxy(t, "/", server.URL), "http://host.example.com/")
	if status := resp.Result().StatusCode; status != 204 {
		t.Errorf("response: want 204, got %d", status)
	}
}

func TestXForwardedProto(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Forwarded-Proto"); got != "https" {
			t.Errorf("X-Forwarded-Proto: want 'https', got %#v", got)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	resp := GET(makeReverseProxy(t, "/", server.URL), "https://example.com/")
	if got := resp.Result().StatusCode; got != 204 {
		t.Errorf("want 204, got %d", got)
	}
}

func TestForwardingForHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "1.1.1.1, "
		if got := r.Header.Get("X-Forwarded-For"); !strings.Contains(got, want) {
			t.Errorf("expect %#v to contain %#v", got, want)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	h := makeReverseProxy(t, "/", server.URL)

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("X-Forwarded-For", "1.1.1.1")
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)
}

func TestForwardingProtoHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Forwarded-Proto"); got != "https" {
			t.Errorf("want 'https', got %#v", got)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	h := makeReverseProxy(t, "/", server.URL)

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)
}

func TestForwardingHostHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "forwarded.example.com"
		if got := r.Host; got != want {
			t.Errorf("Host: want %#v, got %#v", want, got)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	h := makeReverseProxy(t, "/", server.URL)

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("X-Forwarded-Host", "forwarded.example.com")
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)
}

func TestProxyRewrites(t *testing.T) {
	type test struct {
		source, target, request, want string
	}
	tests := []test{
		{"/", "/base", "/api/pets/1", "/base/api/pets/1"},
		{"/api", "/", "/api/pets/1", "/pets/1"},
		{"/api", "/base", "/api/pets/1", "/base/pets/1"},
	}

	for _, tc := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Path; got != tc.want {
				t.Errorf("path: want %#v, got %#v", tc.want, got)
			}
			w.WriteHeader(204)
		}))
		defer server.Close()

		h := makeReverseProxy(t, tc.source, server.URL+tc.target)
		resp := GET(h, tc.request)
		if status := resp.Result().StatusCode; status != 204 {
			t.Errorf("response: want 204, got %d", status)
		}
	}
}

func TestInvalidTargetURL(t *testing.T) {
	h := handler.New()
	err := h.ReverseProxy("/", "://example.com/foo")
	if err == nil {
		t.Error("expect error")
	} else if !strings.Contains(err.Error(), "example.com") {
		t.Errorf("expect %#v to contain 'example.com'", err.Error())
	}
}

func BenchmarkFileServer(b *testing.B) {
	h := handler.New()
	if err := h.FileServer("/", "testdata/html"); err != nil {
		b.Fatalf("/=testdata/html: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := NewBenchmarkResponseWriter()
	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func BenchmarkReverseProxy(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer server.Close()

	h := handler.New()
	if err := h.ReverseProxy("/", server.URL); err != nil {
		b.Fatalf("proxy /=%s: %v", server.URL, err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := NewBenchmarkResponseWriter()
	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, req)
	}
}

func makeFileserver(t *testing.T, source, path string) http.Handler {
	h := handler.New()
	if err := h.FileServer(source, path); err != nil {
		t.Fatalf("fileserver %s=%s: %v", source, path, err)
	}
	return h
}

func makeReverseProxy(t *testing.T, source, URL string) http.Handler {
	h := handler.New()
	if err := h.ReverseProxy(source, URL); err != nil {
		t.Fatalf("proxy %s=%s: %v", source, URL, err)
	}
	return h
}

func GET(h http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", path, nil)
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)

	return resp
}

func readTitle(t *testing.T, resp *http.Response) string {
	body := readBody(t, resp)
	start := strings.Index(body, "<title>")
	end := strings.Index(body, "</title>")
	if start < 0 || end < 0 || start > end {
		t.Fatalf("invalid page:\n%s", body)
	}
	return body[start+7 : end]
}

func readBody(t *testing.T, resp *http.Response) string {
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}
