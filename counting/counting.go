package counting

import (
	"net/http"
	"sync"

	"github.com/mbanzon/middlex"
)

type Counter struct {
	count int64
	mutex *sync.Mutex
}

func New(counters ...*Counter) middlex.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			for _, c := range counters {
				go func(counter *Counter) {
					c.mutex.Lock()
					c.count++
					c.mutex.Unlock()
				}(c)
			}
		})
	}
}

func NewCounter() *Counter {
	c := &Counter{
		mutex: &sync.Mutex{},
	}

	return c
}

func (c *Counter) Count() int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.count
}

func (c *Counter) Reset() {
	c.mutex.Lock()
	c.count = 0
	c.mutex.Unlock()
}
