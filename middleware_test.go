package middlex

import (
	"net/http"
	"testing"
)

type fakeFactory struct {
	counter int
}

func (f *fakeFactory) Middleware() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f.counter++
			h.ServeHTTP(w, r)
		})
	}
}

func TestCombine(t *testing.T) {
	for i := 0; i < 100; i++ {
		f := &fakeFactory{}
		facs := make([]MiddlewareFactory, i)
		for j := 0; j < i; j++ {
			facs[j] = f
		}
		m := Combine(facs...)
		h := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		h.ServeHTTP(nil, nil)

		if f.counter != i {
			t.Fatal("unexpected count:", f.counter)
		}

	}

}
