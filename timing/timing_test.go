package timing

import (
	"math/rand"
	"net/http"
	"sync"
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

func TestTimerWithEmitter(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var totalTime time.Duration

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := time.Duration(rand.Int63n(100)) * time.Millisecond
		totalTime += d
		time.Sleep(d)
	})

	done := false
	wg := sync.WaitGroup{}
	runEmitFn := true

	emitFn := func(ti time.Time, count int64, d time.Duration) {
		if runEmitFn {
			if done {
				if count != 100 {
					t.Fatal("expected count to be 100:", count)
				}
				wg.Done()
			}
		}
	}

	wg.Add(1)

	timer := New(WithEmitter(emitFn, 200*time.Millisecond))
	wrapper := timer.Middleware()
	wrapped := wrapper(handler)

	for i := 0; i < 100; i++ {
		wrapped.ServeHTTP(nil, nil)
	}

	done = true
	wg.Wait()
	runEmitFn = false

	if timer.count != 100 {
		t.Fatal("expected count to be 100:", timer.count)
	}
}

func TestTimerWithEmitReset(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var totalTime time.Duration

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := time.Duration(rand.Int63n(100)) * time.Millisecond
		totalTime += d
		time.Sleep(d)
	})

	done := false
	wg := sync.WaitGroup{}
	var totalCount int64
	runEmitFn := true

	emitFn := func(ti time.Time, count int64, d time.Duration) {
		if runEmitFn {
			totalCount += count
			if done {
				wg.Done()
			}
		}
	}

	wg.Add(1)

	timer := New(WithResetEmitter(emitFn, 200*time.Millisecond))
	wrapper := timer.Middleware()
	wrapped := wrapper(handler)

	for i := 0; i < 100; i++ {
		wrapped.ServeHTTP(nil, nil)
	}

	done = true
	wg.Wait()
	runEmitFn = false

	if timer.count != 0 {
		t.Fatal("expected count to be 0:", timer.count)
	}

	if totalCount != 100 {
		t.Fatal("expected total count to be 100:", totalCount)
	}

}
