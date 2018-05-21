package counting

import (
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestEmitter(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	counts := 0
	signal := make(chan bool)
	var expectedCounts int64

	counter := New(WithEmitter(func(now time.Time, c int64) {
		if counts == 1 {
			<-signal
		}

		if c != expectedCounts {
			t.Fatal("unexpected counts:", c, "vs", expectedCounts)
		}

		counts++
		signal <- false
	}, 500*time.Millisecond))

	wrapper := counter.Middleware()
	wrapped := wrapper(handler)

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

	var totalCount int64
	var expectedCounts int64

	done := false
	signal := make(chan bool)

	counter := New(WithResetEmitter(func(now time.Time, c int64) {
		totalCount += c
		if done {
			if totalCount != expectedCounts {
				t.Fatal("unexpected counts:", totalCount, "vs", expectedCounts)
			}
			signal <- false
		}
	}, 10*time.Millisecond))

	wrapper := counter.Middleware()
	wrapped := wrapper(handler)

	expectedCounts = rand.Int63n(1000) + 1
	for c := int64(0); c < expectedCounts; c++ {
		time.Sleep(time.Millisecond)
		wrapped.ServeHTTP(nil, nil)
	}

	done = true
	<-signal
}
