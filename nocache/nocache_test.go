package nocache

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNocaching(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	noCache := New().Wrap(empty)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	noCache.ServeHTTP(rec, req)

	if rec.Header().Get("Cache-Control") != "no-cache, no-store, must-revalidate" {
		t.Fatal("unexpected Cache-Control header value:", rec.Header().Get("Cache-Control"))
	}

	if rec.Header().Get("Pragma") != "no-cache" {
		t.Fatal("unexpected Pragma header value:", rec.Header().Get("Pragma"))
	}

	if rec.Header().Get("Expires") != "0" {
		t.Fatal("unexpected Expires header value:", rec.Header().Get("Expires"))
	}
}

func TestOptionsRequest(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	noCache := New().Wrap(empty)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodOptions, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	noCache.ServeHTTP(rec, req)

	if rec.Header().Get("Cache-Control") != "" {
		t.Fatal("unexpected Cache-Control header value:", rec.Header().Get("Cache-Control"))
	}

	if rec.Header().Get("Pragma") != "" {
		t.Fatal("unexpected Pragma header value:", rec.Header().Get("Pragma"))
	}

	if rec.Header().Get("Expires") != "" {
		t.Fatal("unexpected Expires header value:", rec.Header().Get("Expires"))
	}
}
