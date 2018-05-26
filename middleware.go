package middlex

import "net/http"

type Middleware func(http.Handler) http.Handler

type MiddlewareFactory interface {
	Middleware() Middleware
}

func Combine(factories ...MiddlewareFactory) Middleware {
	return func(h http.Handler) http.Handler {
		wrapped := h
		for _, f := range factories {
			wrapped = f.Middleware()(wrapped)
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped.ServeHTTP(w, r)
		})
	}
}
