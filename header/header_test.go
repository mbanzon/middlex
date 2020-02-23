package header

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSingleStaticHeader(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hMw := New(WithStaticHeader("X-My-Header", "ValueValueValue"))
	h := hMw.Wrap(emptyHandler)

	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, nil)

	headerValue := recorder.Header().Get("X-My-Header")
	if headerValue != "ValueValueValue" {
		t.Fatal("Unexpected header value:", headerValue)
	}
}

func TestMultipleStaticHeaders(t *testing.T) {
	headers := make(map[string]string)
	headers["X-My-Header1"] = "MyValue1"
	headers["X-My-Header2"] = "MyValue2"
	headers["X-My-Header3"] = "MyValue3"

	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hMw := New(WithStaticHeaders(headers))
	h := hMw.Wrap(emptyHandler)

	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, nil)

	headerVal1 := recorder.Header().Get("X-My-Header1")
	if headerVal1 != "MyValue1" {
		t.Fatal("Unexpected header value for X-My-Header1:", headerVal1)
	}

	headerVal2 := recorder.Header().Get("X-My-Header2")
	if headerVal2 != "MyValue2" {
		t.Fatal("Unexpected header value for X-My-Header2:", headerVal2)
	}

	headerVal3 := recorder.Header().Get("X-My-Header3")
	if headerVal3 != "MyValue3" {
		t.Fatal("Unexpected header value for X-My-Header3:", headerVal3)
	}
}

func TestDynamicHeaderFunc(t *testing.T) {
	hFn := func(r *http.Request) (header, value string) {
		return "MyHeader", "MyValue"
	}

	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hMw := New(WithDynamicHeaderFunc(hFn))
	h := hMw.Wrap(emptyHandler)

	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, nil)

	headerValue := recorder.Header().Get("MyHeader")
	if headerValue != "MyValue" {
		t.Fatal("Unexpected header value:", headerValue)
	}
}

func TestDynamicMultiHeaderFunc(t *testing.T) {
	hFn := func(r *http.Request) (headers map[string]string) {
		tmp := make(map[string]string)
		tmp["MyHeader1"] = "MyValue1"
		tmp["MyHeader2"] = "MyValue2"
		return tmp
	}

	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hMw := New(WithDynamicMultiHeaderFunc(hFn))
	h := hMw.Wrap(emptyHandler)

	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, nil)

	headerValue1 := recorder.Header().Get("MyHeader1")
	if headerValue1 != "MyValue1" {
		t.Fatal("Unexpected header value:", headerValue1)
	}

	headerValue2 := recorder.Header().Get("MyHeader2")
	if headerValue2 != "MyValue2" {
		t.Fatal("Unexpected header value:", headerValue2)
	}
}
