package counting

import (
	"net/http"
	"sync"
)

// Counter allows wrapping of handlers to enable counting of requests.
type Counter struct {
	count int64
	mutex *sync.Mutex
}

type ConfigFunc func(c *Counter)

func (c *Counter) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		c.mutex.Lock()
		c.count++
		c.mutex.Unlock()
	})
}

// New creates a new Counter with the given configuration applied.
func New(configs ...ConfigFunc) *Counter {
	c := &Counter{
		mutex: &sync.Mutex{},
	}

	for _, config := range configs {
		config(c)
	}

	return c
}

// Count returns the current count - the number of requests made through
// the middleware.
func (c *Counter) Count() int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.count
}

// Reset returns the current count (as Count()) but also resets the counter.
func (c *Counter) Reset() int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tmp := c.count
	c.count = 0
	return tmp
}
