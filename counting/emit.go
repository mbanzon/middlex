package counting

import (
	"time"
)

// CountEmitFunc is the signature users of the library should use to provide
// a function used for emitting count information. The function is called
// with the current timestamp and the current count.
type CountEmitFunc func(time.Time, int64)

// WithEmitter returns a ConfigFunc that configures the Counter instance to
// emit the current count by calling the given CountEmitFunc with the given
// time.Duration.
func WithEmitter(fn CountEmitFunc, every time.Duration) ConfigFunc {
	return func(c *Counter) {
		go func() {
			tick := time.Tick(every)
			for now := range tick {
				fn(now, c.Count())
			}
		}()
	}
}

// WithResetEmitter returns a ConfigFunc that configures the Counter instance
// to emit and reset the current count by calling the given CountEmitFunc
// with the given time.Duration.
func WithResetEmitter(fn CountEmitFunc, every time.Duration) ConfigFunc {
	return func(c *Counter) {
		go func() {
			tick := time.Tick(every)
			for now := range tick {
				fn(now, c.Reset())
			}
		}()
	}
}
