package splitter

import (
	"net/http"
	"testing"
)

func TestCreation(t *testing.T) {
	count := 0

	defaultHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
	})

	s := New(WithDefaultHandler(defaultHandler))
	req, _ := http.NewRequest(http.MethodGet, "", nil)

	for i := 0; i < 100; i++ {
		s.ServeHTTP(nil, req)
	}

	if count != 100 {
		t.Fatal()
	}
}

func TestSplit(t *testing.T) {
	defaultCount := 0
	count1 := 0
	count2 := 0

	defaultHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defaultCount++
	})

	handler1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count1++
	})

	handler2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count2++
	})

	s := New(
		WithDefaultHandler(defaultHandler),
		WithPrefix("/api"),
		WithSplit("1", handler1),
		WithSplit("2", handler2),
	)

	req, _ := http.NewRequest(http.MethodGet, "", nil)

	for i := 0; i < 100; i++ {
		s.ServeHTTP(nil, req)
	}

	if defaultCount != 0 && count1 != 0 && count2 != 0 {
		t.Fatal()
	}

	reqDef, _ := http.NewRequest(http.MethodGet, "", nil)
	reqDef.RequestURI = "/api"
	req1, _ := http.NewRequest(http.MethodGet, "", nil)
	req1.RequestURI = "/api/1/bob"
	req2, _ := http.NewRequest(http.MethodGet, "", nil)
	req2.RequestURI = "/api/2"

	for i := 0; i < 100; i++ {
		s.ServeHTTP(nil, reqDef)
		s.ServeHTTP(nil, req1)
		s.ServeHTTP(nil, req2)
	}

	if defaultCount != 100 {
		t.Fatal()
	}

	if count1 != 100 {
		t.Fatal()
	}

	if count2 != 100 {
		t.Fatal()
	}
}
