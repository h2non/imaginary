package main

import "testing"

func TestDefaultError(t *testing.T) {
	err := NewError("oops!\n\n", 503)

	if err.Error() != "oops!" {
		t.Fatal("Invalid error message")
	}
	if err.Code != 503 {
		t.Fatal("Invalid error code")
	}

	code := err.HTTPCode()
	if code != 503 {
		t.Fatalf("Invalid HTTP error status: %d", code)
	}

	json := string(err.JSON())
	if json != "{\"message\":\"oops!\",\"status\":503}" {
		t.Fatalf("Invalid JSON output: %s", json)
	}
}
