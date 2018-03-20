package timing

import (
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestSingleTimer(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := time.Duration(rand.Int63n(100)) * time.Millisecond
		time.Sleep(d)
	})

	timer := NewTimer()
	wrapper := New(timer)
	wrapped := wrapper(handler)

	for i := 0; i < 100; i++ {
		wrapped.ServeHTTP(nil, nil)
	}

	t.Log(timer.Avg())
	t.Log(timer.Count())
}
