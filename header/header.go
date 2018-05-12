package header

import (
	"net/http"

	"github.com/mbanzon/middlex"
)

type Header struct {
	staticHeaders      map[string]string
	dynamicHeaderFuncs []DynamicHeaderFunction
}

type DynamicHeaderFunction func(r *http.Request) (header, value string)

func New(headers map[string]string, headerFuncs ...DynamicHeaderFunction) *Header {
	return &Header{
		staticHeaders:      headers,
		dynamicHeaderFuncs: headerFuncs,
	}
}

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
