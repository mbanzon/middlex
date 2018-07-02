package nocache

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNocaching(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	noCache := New().Middleware()(empty)

	rec := httptest.NewRecorder()
	noCache.ServeHTTP(rec, nil)

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
