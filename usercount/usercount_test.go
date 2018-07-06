package usercount

import (
	"fmt"
	"net/http"
	"testing"
)

func TestUserCountNoResolver(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	counter := New()
	wrapped := counter.Middleware()(empty)

	for i := 0; i < 1000; i++ {
		wrapped.ServeHTTP(nil, nil)
	}

	if counter.GetUserCount() != 1 {
		t.Fatal("unexpected user count:", counter.GetUserCount())
	}

	if counter.GetCount("") != 1000 {
		t.Fatal("unexpected count:", counter.GetCount(""))
	}
}

func TestUserCountIPResolver(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	counter := New(WithIPAddressResolver())
	wrapped := counter.Middleware()(empty)

	for i := 0; i < 1000; i++ {
		for ip := 0; ip < 10; ip++ {
			addr := fmt.Sprintf("10.0.0.%d:12345", ip)
			req, err := http.NewRequest(http.MethodOptions, "/", nil)
			if err != nil {
				t.Fatal("error creating request:", err)
			}
			req.RemoteAddr = addr
			wrapped.ServeHTTP(nil, req)

		}
	}

	if counter.GetUserCount() != 10 {
		t.Fatal("unexpected user count:", counter.GetUserCount())
	}

	for ip := 0; ip < 10; ip++ {
		addr := fmt.Sprintf("10.0.0.%d", ip)
		if counter.GetCount(addr) != 1000 {
			t.Fatal("unexpected count:", counter.GetCount(addr))
		}
	}
}
