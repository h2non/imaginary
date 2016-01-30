package main

import (
	"net/http"
	"net/url"
	"testing"
)

func TestMatchSource(t *testing.T) {
	u, _ := url.Parse("http://foo?url=http://bar/image.jpg")
	req := &http.Request{Method: "GET", URL: u}

	source := MatchSource(req)
	if source == nil {
		t.Error("Cannot match image source")
	}
}
