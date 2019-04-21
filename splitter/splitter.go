package splitter

import (
	"net/http"
	"strings"
)

type Splitter struct {
	prefix         string
	splits         map[string]http.Handler
	defaultHandler http.Handler
}

type ConfigFunc func(*Splitter)

func New(config ...ConfigFunc) *Splitter {
	sh := &Splitter{
		splits: make(map[string]http.Handler),
	}

	for _, confFn := range config {
		confFn(sh)
	}

	return sh
}

func WithPrefix(prefix string) ConfigFunc {
	return func(sh *Splitter) {
		sh.prefix = prefix
	}
}

func WithSplit(path string, h http.Handler) ConfigFunc {
	return func(sh *Splitter) {
		sh.splits[path] = h
	}
}

func WithDefaultHandler(h http.Handler) ConfigFunc {
	return func(sh *Splitter) {
		sh.defaultHandler = h
	}
}

func (sh *Splitter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.RequestURI, sh.prefix) {
		originalURI := r.RequestURI
		r.RequestURI = r.RequestURI[len(sh.prefix):]
		r.URL.Path = r.RequestURI
		path := popURIArg(r)
		hitHandler := false
		for p, h := range sh.splits {
			if path == p {
				hitHandler = true
				h.ServeHTTP(w, r)
				// TODO: should we break here?
			}
		}
		r.RequestURI = originalURI
		r.URL.Path = originalURI
		if !hitHandler && sh.defaultHandler != nil {
			sh.defaultHandler.ServeHTTP(w, r)
		}
	}
}

func popURIArg(r *http.Request) string {
	uri := strings.Trim(r.RequestURI, "/")

	if uri == "" {
		return ""
	}

	var arg string
	index := strings.Index(uri, "/")
	if index != -1 {
		arg = uri[:index]
		r.RequestURI = uri[index:]
	} else {
		arg = uri
		r.RequestURI = ""
	}
	return arg
}
