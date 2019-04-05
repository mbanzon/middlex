package nocache

import (
	"net/http"

	"github.com/mbanzon/middlex/v2/header"
)

type NoCache struct{}

func New() *NoCache {
	nc := &NoCache{}
	return nc
}

func (n *NoCache) Wrap(h http.Handler) http.Handler {
	noCacheHeaders := make(map[string]string)
	noCacheHeaders["Cache-Control"] = "no-cache, no-store, must-revalidate"
	noCacheHeaders["Pragma"] = "no-cache"
	noCacheHeaders["Expires"] = "0"

	return header.New(header.WithDynamicMultiHeaderFunc(func(r *http.Request) (headers map[string]string) {
		if r.Method != http.MethodOptions {
			return noCacheHeaders
		}
		return nil
	})).Wrap(h)
}
