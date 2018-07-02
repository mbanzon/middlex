package nocache

import (
	"github.com/mbanzon/middlex"
	"github.com/mbanzon/middlex/header"
)

type NoCache struct{}

func New() *NoCache {
	return &NoCache{}
}

func (n *NoCache) Middleware() middlex.Middleware {
	noCacheHeaders := make(map[string]string)
	noCacheHeaders["Cache-Control"] = "no-cache, no-store, must-revalidate"
	noCacheHeaders["Pragma"] = "no-cache"
	noCacheHeaders["Expires"] = "0"

	return header.New(header.WithStaticHeaders(noCacheHeaders)).Middleware()
}
