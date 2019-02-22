package timing

import (
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/mbanzon/middlex/v1"
)

var (
	resetValue = time.Duration(math.MaxInt64)
)

type Timer struct {
	count int64
	total time.Duration
	mutex *sync.Mutex
}

func (t *Timer) Middleware() middlex.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			end := time.Now()
			t.mutex.Lock()
			t.count++
			t.total += end.Sub(start)
			t.mutex.Unlock()
		})
	}
}

func New(configs ...ConfigFunc) *Timer {
	t := &Timer{
		mutex: &sync.Mutex{},
	}

	for _, c := range configs {
		c(t)
	}

	return t
}

func (t *Timer) Avg() (int64, time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.count != 0 {
		return t.count, t.total / time.Duration(t.count)
	}

	return t.count, 0
}

func (t *Timer) Reset() (int64, time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	tmpC, tmpAvg := t.count, t.total
	t.count = 0
	t.total = 0
	return tmpC, tmpAvg
}
