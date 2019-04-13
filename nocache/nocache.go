package nocache

import (
	"net/http"
)

type NoCache struct{}

func New() *NoCache {
	nc := &NoCache{}
	return nc
}

func (n *NoCache) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodOptions {
			w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Add("Pragma", "no-cache")
			w.Header().Add("Expires", "0")
		}
	})
}
