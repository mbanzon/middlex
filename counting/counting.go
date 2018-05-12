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

func (c *Counter) Middleware() middlex.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			c.mutex.Lock()
			c.count++
			c.mutex.Unlock()
		})
	}
}

func New() *Counter {
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

func (c *Counter) Reset() int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tmp := c.count
	c.count = 0
	return tmp
}
