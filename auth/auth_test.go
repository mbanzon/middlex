package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAuth(t *testing.T) {
	wrapped, w, r := setupEmptyWrapped(t, http.MethodGet, "/", nil)
	wrapped.ServeHTTP(w, r)
}

func TestOptionsNoAccess(t *testing.T) {
	count := 0
	wrapped := setupWrapped(t, func() {
		count++
	})
	req, _ := http.NewRequest(http.MethodOptions, "", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if count != 0 {
		t.Fatal()
	}

	if w.Code != http.StatusUnauthorized {
		t.Fatal()
	}
}
func TestOptionsAccess(t *testing.T) {
	count := 0
	wrapped := setupWrapped(t, func() {
		count++
	},
		WithOptionsAcceptance())
	req, _ := http.NewRequest(http.MethodOptions, "", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if count != 1 {
		t.Fatal()
	}

	if w.Code != http.StatusOK {
		t.Fatal()
	}
}

func TestWithCookieExtractor(t *testing.T) {
	c := WithCookieTokenExtraction("my-cookie")
	wrapped, w, r := setupEmptyWrapped(t, http.MethodGet, "/", nil, c)
	r.AddCookie(&http.Cookie{Name: "my-cookie", Value: "value"})
	wrapped.ServeHTTP(w, r)
}

func setupEmptyWrapped(t *testing.T, method, url string, data io.Reader, c ...ConfigFunc) (http.Handler, http.ResponseWriter, *http.Request) {
	t.Helper()

	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	auther := New(c...)
	wrapped := auther.Middleware()(empty)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, url, data)

	return wrapped, w, r
}

func setupWrapped(t *testing.T, f func(), c ...ConfigFunc) http.Handler {
	t.Helper()

	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f()
	})
	auther := New(c...)
	wrapped := auther.Middleware()(empty)
	return wrapped
}
