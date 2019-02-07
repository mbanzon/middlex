package usercount

import (
	"time"
)

type CountEmitFunc func(time.Time, int64)

func WithEmitFunction(cef CountEmitFunc, every time.Duration) ConfigFunc {
	return func(u *UserCount) {
		go func() {
			tick := time.Tick(every)
			for now := range tick {
				u.mutex.Lock()
				cef(now, int64(len(u.counts)))
				u.mutex.Unlock()
			}
		}()
	}
}
