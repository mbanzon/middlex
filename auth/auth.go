package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mbanzon/middlex"
)

// ConfigFunc is the type of function used to configure the Authentication
// instance. The library provide various functions that return ConfigFunc
// compatible functions.
type ConfigFunc func(*Authentication)

// LoginFunc is the signature of the function the user if the library
// needs to provide to facilitate the login. The function is called
// when a request is made to the login path (see WithLoginFunc).
type LoginFunc func(r *http.Request) (authorized bool, token string)

// CheckFunc is the function signature for the user-provided function
// used to check the validity of the token. The token is extracted
// from the request and used to call the function on every request.
// The CheckFunc provided should be effective enough to be called
// often.
type CheckFunc func(token string) bool

// LogoutFunc is provided to allow the library to signal that the
// session with the given token should no longer have access.
type LogoutFunc func(token string)

type tokenExtractFunc func(*http.Request) (found bool, token string)

// Authentication holds the functions and data configured and provide
// the middleware used for authentication.
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

// New constructs a new Authentication instance and applies the configuration
// given. The default token extraction method is the "Authorization" header
// (expected to hold a "Bearer " prefix).
func New(config ...ConfigFunc) *Authentication {
	a := &Authentication{
		tokenExFn: extractBearerToken,
	}
	for _, confFn := range config {
		confFn(a)
	}

	return a
}

// WithCookieTokenExtraction creates a ConfigFunc that configures the
// Authentication to extract the authorization token from a cookie with
// the given name.
func WithCookieTokenExtraction(name string) ConfigFunc {
	return func(a *Authentication) {
		a.tokenExFn = func(r *http.Request) (bool, string) {
			return extractCookieToken(name, r)
		}
	}
}

// WithHeaderTokenExtraction creates as ConfigFunc that configures the
// Authentication to extract the authorization token from a header with
// the given name.
func WithHeaderTokenExtraction(name string) ConfigFunc {
	return func(a *Authentication) {
		a.tokenExFn = func(r *http.Request) (bool, string) {
			return extractHeaderToken(name, r)
		}
	}
}

// WithAuthurizationHeaderTokenExtraction creates as ConfigFunc that
// configures the Authentication to extract the authorization token
// from the Authorization header.
func WithAuthurizationHeaderTokenExtraction() ConfigFunc {
	return WithHeaderTokenExtraction("Authorization")
}

// WithCheckFunc creates a ConfigFunc that configures an Authentication
// instance to use the given CheckFunc.
func WithCheckFunc(checkFn CheckFunc) ConfigFunc {
	return func(a *Authentication) {
		a.checkFn = checkFn
	}
}

// WithLoginFunc creates a ConfigFunc that configures an Authentication
// instance to use the given LoginFunc and call it when the given loginPath
// is requested.
func WithLoginFunc(loginPath string, loginFn LoginFunc) ConfigFunc {
	return func(a *Authentication) {
		a.loginPath = loginPath
		a.loginFn = loginFn
	}
}

// WithLogoutFunc creates a ConfigFunc that configures an Authentication
// instance to call the given LogoutFunc when the given logoutPath is
// requested.
func WithLogoutFunc(logoutPath string, logoutFn LogoutFunc) ConfigFunc {
	return func(a *Authentication) {
		a.logoutPath = logoutPath
		a.logoutFn = logoutFn
	}
}

// WithExcludedPaths creates a ConfigFunc that configures an Authentication
// to bypass the paths given when they are requested.
func WithExcludedPaths(paths ...string) ConfigFunc {
	return func(a *Authentication) {
		for _, path := range paths {
			a.excludedPaths = append(a.excludedPaths, path)
		}
	}
}

// WithExcludedPrefixes creates a ConfigFunc that configures an
// Authentication to bypass paths with one of the given prefixes
// when they are requested.
func WithExcludedPrefixes(prefixes ...string) ConfigFunc {
	return func(a *Authentication) {
		for _, prefix := range prefixes {
			a.excludedPrefixes = append(a.excludedPrefixes, prefix)
		}
	}
}

// Middleware retuns the middlex.Middleware for the Authentication.
// The returned value can be used to wrap a http.Handler in this
// Authentication.
func (a *Authentication) Middleware() middlex.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if a.loginFn != nil && r.RequestURI == a.loginPath && r.Method == http.MethodPost {
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
				if a.logoutFn != nil && r.RequestURI == a.logoutPath && r.Method == http.MethodGet {
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
		return true, strings.TrimPrefix(headerStr, "Bearer ")
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

func extractHeaderToken(name string, r *http.Request) (found bool, token string) {
	token = r.Header.Get(name)
	return token != "", token
}
