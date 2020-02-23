package usercount

import (
	"fmt"
	"net/http"
	"testing"
)

func TestUserCountNoResolver(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	counter := New()
	wrapped := counter.Wrap(empty)

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
	wrapped := counter.Wrap(empty)

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

		counter.Reset(addr)
		if counter.GetCount(addr) != 0 {
			t.Fatal("unexpected count after reset:", counter.GetCount(addr))
		}
	}
}

func TestUserCountAuthenticationResolver(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	counter := New(WithAuthenticationResolver())
	wrapped := counter.Wrap(empty)

	for i := 0; i < 1000; i++ {
		for ip := 0; ip < 10; ip++ {
			bearer := fmt.Sprintf("Bearer %d", ip)
			req, err := http.NewRequest(http.MethodOptions, "/", nil)
			if err != nil {
				t.Fatal("error creating request:", err)
			}
			req.Header.Set("Authorization", bearer)
			wrapped.ServeHTTP(nil, req)

		}
	}

	if counter.GetUserCount() != 10 {
		t.Fatal("unexpected user count:", counter.GetUserCount())
	}

	for ip := 0; ip < 10; ip++ {
		addr := fmt.Sprintf("Bearer %d", ip)
		if counter.GetCount(addr) != 1000 {
			t.Fatal("unexpected count:", counter.GetCount(addr))
		}
	}

	counter.ResetAll()

	for ip := 0; ip < 10; ip++ {
		addr := fmt.Sprintf("Bearer %d", ip)
		if counter.GetCount(addr) != 0 {
			t.Fatal("unexpected count after resetting all:", counter.GetCount(addr))
		}
	}
}

func TestUserCountFailingIPResolver(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	counter := New(WithIPAddressResolver())
	wrapped := counter.Wrap(empty)

	for i := 0; i < 1000; i++ {
		req, err := http.NewRequest(http.MethodOptions, "/", nil)
		if err != nil {
			t.Fatal("error creating request:", err)
		}
		wrapped.ServeHTTP(nil, req)
	}

	if counter.GetUserCount() != 1 {
		t.Fatal("unexpected user count:", counter.GetUserCount())
	}

	if counter.GetCount("") != 1000 {
		t.Fatal("unexpected count:", counter.GetCount(""))
	}
}

func TestUserCountCookieResolver(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	counter := New(WithCookieResolver("UserCookie"))
	wrapped := counter.Wrap(empty)

	for i := 0; i < 1000; i++ {
		for c := 0; c < 10; c++ {
			cookie := &http.Cookie{Name: "UserCookie", Value: fmt.Sprintf("%d", c)}
			req, err := http.NewRequest(http.MethodOptions, "/", nil)
			if err != nil {
				t.Fatal("error creating request:", err)
			}
			req.AddCookie(cookie)
			wrapped.ServeHTTP(nil, req)

		}
	}

	if counter.GetUserCount() != 10 {
		t.Fatal("unexpected user count:", counter.GetUserCount())
	}

	for c := 0; c < 10; c++ {
		cookieValue := fmt.Sprintf("%d", c)
		if counter.GetCount(cookieValue) != 1000 {
			t.Fatal("unexpected count:", counter.GetCount(cookieValue))
		}
	}
}

func TestUserCountFailingCookieResolver(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	counter := New(WithCookieResolver("UserCookie"))
	wrapped := counter.Wrap(empty)

	for i := 0; i < 1000; i++ {
		req, err := http.NewRequest(http.MethodOptions, "/", nil)
		if err != nil {
			t.Fatal("error creating request:", err)
		}
		wrapped.ServeHTTP(nil, req)
	}

	if counter.GetUserCount() != 1 {
		t.Fatal("unexpected user count:", counter.GetUserCount())
	}

	if counter.GetCount("") != 1000 {
		t.Fatal("unexpected count:", counter.GetCount(""))
	}
}

func TestUserCountCustomResolver(t *testing.T) {
	empty := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	resolver := func(*http.Request) string {
		return ""
	}
	counter := New(WithCustomResolver(resolver))
	wrapped := counter.Wrap(empty)

	for i := 0; i < 1000; i++ {
		req, err := http.NewRequest(http.MethodOptions, "/", nil)
		if err != nil {
			t.Fatal("error creating request:", err)
		}
		wrapped.ServeHTTP(nil, req)
	}

	if counter.GetUserCount() != 1 {
		t.Fatal("unexpected user count:", counter.GetUserCount())
	}

	if counter.GetCount("") != 1000 {
		t.Fatal("unexpected count:", counter.GetCount(""))
	}
}
