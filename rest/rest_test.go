package rest

import "testing"

func TestNew(t *testing.T) {
	name := "test"

	tmp := New(name)
	if tmp == nil {
		t.Fatal("expected nil:", tmp)
	}

	if name != tmp.resourceName {
		t.Fatalf("expected name to be: %s, but found: %s", name, tmp.resourceName)
	}
}
