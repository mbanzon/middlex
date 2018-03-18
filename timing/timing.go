package timing

import (
	"net/http"
	"sync"
	"time"

	"github.com/mbanzon/middlex"
)

type Timer struct {
	count int64
	total time.Duration
	lock  sync.Mutex
}

func New(timers ...*Timer) middlex.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			end := time.Now()
			go func() {
				for _, timer := range timers {
					d := end.Sub(start)

					timer.lock.Lock()
					timer.count++
					timer.total += d
					timer.lock.Unlock()
				}
			}()
		})
	}
}

func (t *Timer) Avg() time.Duration {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.count == 0 {
		return 0
	}

	return t.total / time.Duration(t.count)
}

func (t *Timer) Count() int64 {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.count
}

func (t *Timer) Reset() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.count = 0
	t.total = 0
}
