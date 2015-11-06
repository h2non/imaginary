package main

import "testing"

func TestError(t *testing.T) {
	err := NewError("oops!\n\n", 1)

	if err.Error() != "oops!" {
		t.Fatal("Invalid error message")
	}
	if err.Code != 1 {
		t.Fatal("Invalid error code")
	}

	code := err.HTTPCode()
	if code != 400 {
		t.Fatalf("Invalid HTTP error status: %d", code)
	}

	json := string(err.JSON())
	if json != "{\"message\":\"oops!\",\"code\":1}" {
		t.Fatalf("Invalid JSON output: %s", json)
	}
}
