package counting

import (
	"time"
)

type CountEmitFunc func(time.Time, int64)

func (c *Counter) Emit(every time.Duration, receiver CountEmitFunc) {
	tick := time.Tick(every)
	for now := range tick {
		receiver(now, c.Count())
	}
}

func (c *Counter) EmitReset(every time.Duration, receiver CountEmitFunc) {
	tick := time.Tick(every)
	for now := range tick {
		receiver(now, c.Reset())
	}
}
