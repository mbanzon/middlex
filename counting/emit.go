package counting

import (
	"time"
)

// CountEmitFunc is the signature users of the library should use to provide
// a function used for emitting count information. The function is called
// with the current timestamp and the current count.
type CountEmitFunc func(time.Time, int64)

// Emit contructs and starts a new go routine that continously emits the
// current count of the counter by calling the provided receiver function
// in the interval given.
func (c *Counter) Emit(every time.Duration, receiver CountEmitFunc) {
	go func() {
		tick := time.Tick(every)
		for now := range tick {
			receiver(now, c.Count())
		}
	}()
}

// EmitReset contructs and starts a new go routine that continously emits and
// resets the current count of the counter by calling the provided receiver
// function in the interval given.
func (c *Counter) EmitReset(every time.Duration, receiver CountEmitFunc) {
	go func() {
		tick := time.Tick(every)
		for now := range tick {
			receiver(now, c.Reset())
		}
	}()
}
