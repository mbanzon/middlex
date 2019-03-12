package usercount

import (
	"net"
	"net/http"
	"sync"

	"github.com/mbanzon/middlex/v2"
)

type UserCount struct {
	mutex           *sync.Mutex
	counts          map[string]int64
	userResolveFunc ResolverFunc
}

type ResolverFunc func(*http.Request) string

type ConfigFunc func(*UserCount)

func New(configs ...ConfigFunc) *UserCount {
	u := &UserCount{
		mutex:           &sync.Mutex{},
		counts:          make(map[string]int64),
		userResolveFunc: func(*http.Request) string { return "" },
	}

	for _, c := range configs {
		c(u)
	}

	return u
}

func (u *UserCount) GetUserCount() int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return int64(len(u.counts))
}

func (u *UserCount) GetCount(user string) int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.counts[user]
}

func (u *UserCount) Reset(user string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.counts[user] = 0
}

func (u *UserCount) ResetAll() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.counts = make(map[string]int64)
}

func (u *UserCount) Middleware() middlex.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := u.userResolveFunc(r)
			u.mutex.Lock()
			count := u.counts[user]
			count++
			u.counts[user] = count
			u.mutex.Unlock()

			h.ServeHTTP(w, r)
		})
	}
}

func WithCustomResolver(fn ResolverFunc) ConfigFunc {
	return func(u *UserCount) {
		u.userResolveFunc = fn
	}
}

func WithAuthenticationResolver() ConfigFunc {
	return WithHeaderResolver("Authorization")
}

func WithHeaderResolver(header string) ConfigFunc {
	return func(u *UserCount) {
		u.userResolveFunc = func(r *http.Request) string {
			return r.Header.Get(header)
		}
	}
}

func WithCookieResolver(name string) ConfigFunc {
	return func(u *UserCount) {
		u.userResolveFunc = func(r *http.Request) string {
			c, err := r.Cookie(name)
			if err != nil {
				return ""
			}

			return c.Value
		}
	}
}

func WithIPAddressResolver() ConfigFunc {
	return func(u *UserCount) {
		u.userResolveFunc = func(r *http.Request) string {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				return ""
			}
			return host
		}
	}
}
