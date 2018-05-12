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

	timer := New()
	wrapper := timer.Middleware()
	wrapped := wrapper(handler)

	t.Log(timer.Avg())

	for i := 0; i < 100; i++ {
		wrapped.ServeHTTP(nil, nil)
	}

	timer.Reset()
}
