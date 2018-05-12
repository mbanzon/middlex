package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mbanzon/middlex"
)

type ConfigFunc func(*Authentication)

type LoginFunc func(r *http.Request) (authorized bool, token string)

type CheckFunc func(token string) bool

type LogoutFunc func(token string)

type tokenExtractFunc func(*http.Request) (found bool, token string)

type Authentication struct {
	loginPath        string
	loginFn          LoginFunc
	logoutPath       string
	logoutFn         LogoutFunc
	checkFn          CheckFunc
	excludedPaths    []string
	excludedPrefixes []string
	tokenExFn        tokenExtractFunc
}

func New(config ...ConfigFunc) *Authentication {
	a := &Authentication{
		tokenExFn: extractBearerToken,
	}
	for _, confFn := range config {
		confFn(a)
	}

	return a
}

func WithCookie(name string) ConfigFunc {
	return func(a *Authentication) {
		a.tokenExFn = func(r *http.Request) (bool, string) {
			return extractCookieToken(name, r)
		}
	}
}

func WithCheckFunc(checkFn CheckFunc) ConfigFunc {
	return func(a *Authentication) {
		a.checkFn = checkFn
	}
}

func WithLoginFunc(loginPath string, loginFn LoginFunc) ConfigFunc {
	return func(a *Authentication) {
		a.loginPath = loginPath
		a.loginFn = loginFn
	}
}

func WithLogoutFunc(logoutPath string, logoutFn LogoutFunc) ConfigFunc {
	return func(a *Authentication) {
		a.logoutPath = logoutPath
		a.logoutFn = logoutFn
	}
}

func WithExcludedPaths(paths ...string) ConfigFunc {
	return func(a *Authentication) {
		for _, path := range paths {
			a.excludedPaths = append(a.excludedPaths, path)
		}
	}
}

func WithExcludedPrefixes(prefixes ...string) ConfigFunc {
	return func(a *Authentication) {
		for _, prefix := range prefixes {
			a.excludedPrefixes = append(a.excludedPrefixes, prefix)
		}
	}
}

func (a *Authentication) Middleware() middlex.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if a.loginFn != nil && r.RequestURI == a.loginPath {
				authorized, token := a.loginFn(r)
				if !authorized {
					http.Error(w, "authorization failed", http.StatusUnauthorized)
					return
				}
				fmt.Fprint(w, token)
				return
			}

			for _, ex := range a.excludedPaths {
				if ex == r.RequestURI {
					h.ServeHTTP(w, r)
					return
				}
			}

			for _, exp := range a.excludedPrefixes {
				if strings.HasPrefix(r.RequestURI, exp) {
					h.ServeHTTP(w, r)
					return
				}
			}

			hasAuthHeader, token := a.tokenExFn(r)

			if hasAuthHeader {
				if a.logoutFn != nil && r.RequestURI == a.logoutPath {
					a.logoutFn(token)
					return
				}

				if a.checkFn != nil && a.checkFn(token) {
					h.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "", http.StatusUnauthorized)
		})
	}
}

func extractBearerToken(r *http.Request) (found bool, token string) {
	headerStr := r.Header.Get("Authorization")
	if headerStr != "" {
		return true, strings.TrimPrefix(headerStr, "Bearer")
	}
	return false, ""
}

func extractCookieToken(name string, r *http.Request) (found bool, token string) {
	c, err := r.Cookie(name)
	if err != nil {
		return false, ""
	}

	return true, c.Value
}
