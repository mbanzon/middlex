package header

import (
	"net/http"

	"github.com/mbanzon/middlex"
)

// Header allows wrapping of handlers to provide a structured way of adding
// headers to HTTP responses. Header can handle two types of headers - static
// headers and dynamic headers. Static headers are fixed value and dynamic
// headers are resolved at call time.
type Header struct {
	staticHeaders      map[string]string
	dynamicHeaderFuncs []DynamicHeaderFunction
}

// DynamicHeaderFunction is the signature of the functions that can be given
// to resolve header values dynamically.
type DynamicHeaderFunction func(r *http.Request) (header, value string)

// New creates a new Header with the given static headers and the given
// dynamic header functions.
func New(headers map[string]string, headerFuncs ...DynamicHeaderFunction) *Header {
	return &Header{
		staticHeaders:      headers,
		dynamicHeaderFuncs: headerFuncs,
	}
}

// Middleware returns the middlex.Middleware for the Header that can be
// used to wrap handlers.
func (hm *Header) Middleware() middlex.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for header, value := range hm.staticHeaders {
				w.Header().Add(header, value)
			}
			for _, hFn := range hm.dynamicHeaderFuncs {
				header, value := hFn(r)
				w.Header().Add(header, value)
			}
			h.ServeHTTP(w, r)
		})
	}
}
