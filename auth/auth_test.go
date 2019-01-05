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
