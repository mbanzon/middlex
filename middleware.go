package middlex

import "net/http"

// Middleware describes the function used to wrap a http.Handler in another
// to apply middleware.
type Middleware func(http.Handler) http.Handler

// MiddlewareFactory describes the interface used to create Middleware
// instances.
type MiddlewareFactory interface {
	Middleware() Middleware
}

// Combine takes a number of MiddlewareFactory and use them to combine them
// into a single Middleware. Please note that on every invokation of the
// resulting Middleware the factories are used to produce new Middleware
// instances.
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
