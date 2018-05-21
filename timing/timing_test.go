package timing

import (
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestSingleTimer(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var totalTime time.Duration

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := time.Duration(rand.Int63n(100)) * time.Millisecond
		totalTime += d
		time.Sleep(d)
	})

	timer := New()
	wrapper := timer.Middleware()
	wrapped := wrapper(handler)

	for i := 0; i < 100; i++ {
		wrapped.ServeHTTP(nil, nil)
	}

	count, avg := timer.Avg()
	if count != 100 {
		t.Fatal("unexptected count:", count)
	}

	if avg < totalTime/100 {
		t.Fatal("avg too low:", avg)
	}

	timer.Reset()

	count, avg = timer.Avg()

	if count != 0 {
		t.Fatal("unexpected count after reset:", count)
	}

	if avg != 0 {
		t.Fatal("unexpected avg after reset:", avg)
	}
}
