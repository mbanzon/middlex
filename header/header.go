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

type ConfigFunc func(*Header)

// DynamicHeaderFunction is the signature of the functions that can be given
// to resolve header values dynamically.
type DynamicHeaderFunction func(r *http.Request) (header, value string)

// New creates a new Header with the given static headers and the given
// dynamic header functions.
func New(config ...ConfigFunc) *Header {
	h := &Header{
		staticHeaders:      make(map[string]string),
		dynamicHeaderFuncs: nil,
	}
	for _, c := range config {
		c(h)
	}
	return h
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

func WithDynamicHeaderFunc(dFn DynamicHeaderFunction) ConfigFunc {
	return func(h *Header) {
		h.dynamicHeaderFuncs = append(h.dynamicHeaderFuncs, dFn)
	}
}

func WithStaticHeader(header, value string) ConfigFunc {
	return func(h *Header) {
		h.staticHeaders[header] = value
	}
}

func WithStaticHeaders(headers map[string]string) ConfigFunc {
	return func(h *Header) {
		for header, value := range headers {
			h.staticHeaders[header] = value
		}
	}
}
