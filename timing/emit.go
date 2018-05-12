package timing

import (
	"time"
)

type ConfigFunc func(*Timer)

type AvgTimeEmitFunc func(time.Time, int64, time.Duration)

func WithEmitter(fn AvgTimeEmitFunc, every time.Duration) ConfigFunc {
	return func(t *Timer) {
		go func() {
			ticker := time.Tick(every)
			for now := range ticker {
				tmpC, tmpAvg := t.Avg()
				fn(now, tmpC, tmpAvg)
			}
		}()
	}
}

func WithResetEmitter(fn AvgTimeEmitFunc, every time.Duration) ConfigFunc {
	return func(t *Timer) {
		go func() {
			ticker := time.Tick(every)
			for now := range ticker {
				tmpC, tmpAvg := t.Reset()
				fn(now, tmpC, tmpAvg)
			}
		}()
	}
}
