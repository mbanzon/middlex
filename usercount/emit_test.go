package usercount

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestUserCountEmitting(t *testing.T) {
	var emittedCount int64
	emitFn := func(ts time.Time, count int64) {
		fmt.Println("in here")
		emittedCount += count
	}

	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	counter := New(WithEmitFunction(emitFn, 100*time.Millisecond), WithHeaderResolver("X-User"))
	wrapped := counter.Wrap(empty)

	for u := 0; u < 10; u++ {
		for i := 0; i < 1000; i++ {
			req, err := http.NewRequest(http.MethodOptions, "/", nil)
			if err != nil {
				t.Fatal("error creating request:", err)
			}
			req.Header.Set("X-User", fmt.Sprint(u))
			wrapped.ServeHTTP(nil, req)
		}
	}

	done := make(chan bool)
	go func() {
		for {
			if emittedCount == 10 {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		done <- true
	}()

	select {
	case <-done:
		break
	case <-time.After(5 * time.Second):
		t.Fatal("waited too long")
	}

	if counter.GetUserCount() != 10 {
		t.Fatal("unexpected user count:", counter.GetUserCount())
	}

	for u := 0; u < 10; u++ {
		user := fmt.Sprint(u)
		if counter.GetCount(user) != 1000 {
			t.Fatal("unexpected count:", counter.GetCount(user))
		}
	}
}
