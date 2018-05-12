package cors

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mbanzon/middlex"
)

type Cors struct {
	allowedOrigins []string
	allowedHeaders []string
	allowedMethods []string
	maxAge         time.Duration
}

type ConfigFunc func(*Cors)

func New(configs ...ConfigFunc) *Cors {
	c := &Cors{}

	for _, cFn := range configs {
		cFn(c)
	}

	return c
}

func (c *Cors) Middleware() middlex.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(c.allowedOrigins) > 0 {
				w.Header().Add("Access-Control-Allow-Origin", strings.Join(c.allowedOrigins, ", "))
			}
			if len(c.allowedMethods) > 0 {
				w.Header().Add("Access-Control-Allow-Methods", strings.Join(c.allowedMethods, ", "))
			}
			if len(c.allowedHeaders) > 0 {
				w.Header().Add("Access-Control-Allow-Headers", strings.Join(c.allowedHeaders, ", "))
			}

			if r.Method == http.MethodOptions {
				if c.maxAge > 0 {
					w.Header().Add("Access-Control-Max-Age", fmt.Sprint(int(c.maxAge.Seconds())))
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func WithOrigins(origins ...string) ConfigFunc {
	return func(c *Cors) {
		c.allowedOrigins = origins
	}
}

func WithMethods(methods ...string) ConfigFunc {
	return func(c *Cors) {
		c.allowedMethods = methods
	}
}

func WithMaxAge(age time.Duration) ConfigFunc {
	return func(c *Cors) {
		c.maxAge = age
	}
}

func WithHeaders(headers ...string) ConfigFunc {
	return func(c *Cors) {
		c.allowedHeaders = headers
	}
}
