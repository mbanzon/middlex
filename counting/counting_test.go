package counting

import (
	"net/http"
	"testing"
	"time"
)

func TestSingleCounter(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	counter := NewCounter()
	wrapper := New(counter)
	wrapped := wrapper(handler)

	for i := 0; i < 10; i++ {
		// count := rand.Int63n(10)
		// t.Log(count)
		// t.Log(counter.Count())
		// for c := int64(0); c < count; i++ {
		wrapped.ServeHTTP(nil, nil)

		// }

	}

	time.Sleep(time.Second)

	if 10 != counter.Count() {
		t.Fatal("unexpected count:", 10, "vs", counter.Count())
	}

	counter.Reset()
	if counter.Count() != 0 {
		t.Fatal("expected count to be zero:", counter.Count())
	}

}
