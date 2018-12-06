package nocache

import (
	"net/http"

	"github.com/mbanzon/middlex"
	"github.com/mbanzon/middlex/header"
)

type ConfigFunc func(*NoCache)

type NoCache struct {
}

func New(config ...ConfigFunc) *NoCache {
	nc := &NoCache{}
	for _, c := range config {
		c(nc)
	}
	return nc
}

func (n *NoCache) Middleware() middlex.Middleware {
	noCacheHeaders := make(map[string]string)
	noCacheHeaders["Cache-Control"] = "no-cache, no-store, must-revalidate"
	noCacheHeaders["Pragma"] = "no-cache"
	noCacheHeaders["Expires"] = "0"

	return header.New(header.WithDynamicMultiHeaderFunc(func(r *http.Request) (headers map[string]string) {
		if r.Method != http.MethodOptions {
			return noCacheHeaders
		}
		return nil
	})).Middleware()
}
