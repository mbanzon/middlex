package healthcheck

import (
	"net/http"
	"time"
)

type HealthChecker struct {
	responseCode  int
	message       string
	checker       CheckerFunc
	checkInterval time.Duration
}

type CheckerFunc func() bool

type ConfigFunc func(*HealthChecker)

func New(c ...ConfigFunc) *HealthChecker {
	h := &HealthChecker{
		responseCode: http.StatusInternalServerError,
	}

	for _, cf := range c {
		cf(h)
	}

	return h
}

func (hc *HealthChecker) Wrap(h http.Handler) http.Handler {
	if hc.checker != nil {
		if hc.checkInterval > 0 {
			healthy := true
			go func() {
				for {
					healthy = hc.checker()
					time.Sleep(hc.checkInterval)
				}
			}()
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !healthy {
					http.Error(w, hc.message, hc.responseCode)
					return
				}
				h.ServeHTTP(w, r)
			})
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !hc.checker() {
				http.Error(w, hc.message, hc.responseCode)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func WithResponseCode(code int) ConfigFunc {
	return func(h *HealthChecker) {
		h.responseCode = code
	}
}

func WithMessage(msg string) ConfigFunc {
	return func(h *HealthChecker) {
		h.message = msg
	}
}

func WithCheckerFunc(chk CheckerFunc) ConfigFunc {
	return func(h *HealthChecker) {
		h.checker = chk
	}
}

func WithInterval(d time.Duration) ConfigFunc {
	return func(h *HealthChecker) {
		h.checkInterval = d
	}
}
