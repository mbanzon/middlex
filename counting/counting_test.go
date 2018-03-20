package counting

import (
	"math/rand"
	"net/http"
	"testing"
)

func TestSingleCounter(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	counter := NewCounter()
	wrapper := New(counter)
	wrapped := wrapper(handler)

	for i := 0; i < 10; i++ {
		count := rand.Int63n(10)
		for c := int64(0); c < count; c++ {
			wrapped.ServeHTTP(nil, nil)
		}

		if count != counter.Count() {
			t.Fatal("unexpected count:", count, "vs", counter.Count())
		}

		counter.Reset()
		if counter.Count() != 0 {
			t.Fatal("expected count to be zero:", counter.Count())
		}
	}
}
