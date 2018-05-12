package counting

import (
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestEmitter(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	counter := New()
	wrapper := counter.Middleware()
	wrapped := wrapper(handler)

	counts := 0
	signal := make(chan bool)
	var expectedCounts int64

	counter.Emit(500*time.Millisecond, func(now time.Time, c int64) {
		if counts == 1 {
			<-signal
		}

		if c != expectedCounts {
			t.Fatal("unexpected counts:", c, "vs", expectedCounts)
		}

		counts++
		signal <- false
	})

	<-signal

	expectedCounts = rand.Int63n(1000) + 1
	for c := int64(0); c < expectedCounts; c++ {
		wrapped.ServeHTTP(nil, nil)
	}

	signal <- false
	<-signal
}

func TestResetEmitter(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	counter := New()
	wrapper := counter.Middleware()
	wrapped := wrapper(handler)

	done := false

	var totalCount int64
	var expectedCounts int64

	signal := make(chan bool)

	counter.EmitReset(10*time.Millisecond, func(now time.Time, c int64) {
		totalCount += c
		if done {
			if totalCount != expectedCounts {
				t.Fatal("unexpected counts:", totalCount, "vs", expectedCounts)
			}
			signal <- false
		}
	})

	expectedCounts = rand.Int63n(1000) + 1
	for c := int64(0); c < expectedCounts; c++ {
		time.Sleep(time.Millisecond)
		wrapped.ServeHTTP(nil, nil)
	}

	done = true
	<-signal
}
