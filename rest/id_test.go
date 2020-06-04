package rest

import "testing"

func TestSimpleIDParsing(t *testing.T) {
	sampleURI := "/some/resource/42"
	id := ParseRESTResourceID(sampleURI, "resource")

	if !id.Present() {
		t.Fatal("id should be present")
	}

	if id.GetError() != nil {
		t.Fatal("id should be valid int")
	}

	if id.Get() != 42 {
		t.Fatal("id should be int64=42")
	}

	if id.GetAsStringError() != nil {
		t.Fatal("id should be valid string")
	}

	if id.GetAsString() != "42" {
		t.Fatal("id should be string=\"42\"")
	}
}

func TestSimpleIDParsingFail(t *testing.T) {
	sampleURI := "/some/resource"
	id := ParseRESTResourceID(sampleURI, "resource")

	if id.Present() {
		t.Fatal("id should not be present")
	}

	if id.GetError() == nil {
		t.Fatal("id should not be valid int")
	}

	if id.GetAsStringError() == nil {
		t.Fatal("id should not be valid string")
	}
}
