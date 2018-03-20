package timing

import (
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/mbanzon/middlex"
)

var (
	resetValue = time.Duration(math.MaxInt64)
)

type Timer struct {
	count int64
	total time.Duration
	mutex *sync.Mutex
}

func New(timers ...*Timer) middlex.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			end := time.Now()
			for _, t := range timers {
				go func(timer *Timer) {
					timer.mutex.Lock()
					timer.count++
					timer.total += end.Sub(start)
					timer.mutex.Unlock()
				}(t)
			}
		})
	}
}

func NewTimer() *Timer {
	t := &Timer{
		mutex: &sync.Mutex{},
	}

	return t
}

func (t *Timer) Avg() time.Duration {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.count != 0 {
		return t.total / time.Duration(t.count)
	}

	return 0
}

func (t *Timer) Count() int64 {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.count
}

func (t *Timer) Reset() {
	t.mutex.Lock()
	t.count = 0
	t.total = 0
	t.mutex.Unlock()
}
