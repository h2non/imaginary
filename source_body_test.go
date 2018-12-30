package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

const fixtureFile = "testdata/large.jpg"

func TestSourceBodyMatch(t *testing.T) {
	u, _ := url.Parse("http://foo")
	req := &http.Request{Method: http.MethodPost, URL: u}
	source := NewBodyImageSource(&SourceConfig{})

	if !source.Matches(req) {
		t.Error("Cannot match the request")
	}
}

func TestBodyImageSource(t *testing.T) {
	var body []byte
	var err error

	source := NewBodyImageSource(&SourceConfig{})
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		body, err = source.GetImage(r)
		if err != nil {
			t.Fatalf("Error while reading the body: %s", err)
		}
		w.Write(body)
	}

	file, _ := os.Open(fixtureFile)
	r, _ := http.NewRequest(http.MethodPost, "http://foo/bar", file)
	w := httptest.NewRecorder()
	fakeHandler(w, r)

	buf, _ := ioutil.ReadFile(fixtureFile)
	if len(body) != len(buf) {
		t.Error("Invalid response body")
	}
}

func testReadBody(t *testing.T) {
	var body []byte
	var err error

	source := NewBodyImageSource(&SourceConfig{})
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		body, err = source.GetImage(r)
		if err != nil {
			t.Fatalf("Error while reading the body: %s", err)
		}
		w.Write(body)
	}

	file, _ := os.Open(fixtureFile)
	r, _ := http.NewRequest(http.MethodPost, "http://foo/bar", file)
	w := httptest.NewRecorder()
	fakeHandler(w, r)

	buf, _ := ioutil.ReadFile(fixtureFile)
	if len(body) != len(buf) {
		t.Error("Invalid response body")
	}
}
