package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const fixtureImage = "testdata/large.jpg"
const fixture1024Bytes = "testdata/1024bytes"

func TestHttpImageSource(t *testing.T) {
	var body []byte
	var err error

	buf, _ := ioutil.ReadFile(fixtureImage)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(buf)
	}))
	defer ts.Close()

	source := NewHttpImageSource(&SourceConfig{})
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

	r, _ := http.NewRequest("GET", "http://foo/bar?url="+ts.URL, nil)
	w := httptest.NewRecorder()
	fakeHandler(w, r)

	if len(body) != len(buf) {
		t.Error("Invalid response body")
	}
}

func TestHttpImageSourceAllowedOrigin(t *testing.T) {
	buf, _ := ioutil.ReadFile(fixtureImage)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(buf)
	}))
	defer ts.Close()

	origin, _ := url.Parse(ts.URL)
	origins := []*url.URL{origin}
	source := NewHttpImageSource(&SourceConfig{AllowedOrigings: origins})

	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		body, err := source.GetImage(r)
		if err != nil {
			t.Fatalf("Error while reading the body: %s", err)
		}
		w.Write(body)

		if len(body) != len(buf) {
			t.Error("Invalid response body length")
		}
	}

	r, _ := http.NewRequest("GET", "http://foo/bar?url="+ts.URL, nil)
	w := httptest.NewRecorder()
	fakeHandler(w, r)
}

func TestHttpImageSourceNotAllowedOrigin(t *testing.T) {
	origin, _ := url.Parse("http://foo")
	origins := []*url.URL{origin}
	source := NewHttpImageSource(&SourceConfig{AllowedOrigings: origins})

	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		_, err := source.GetImage(r)
		if err == nil {
			t.Fatal("Error cannot be empty")
		}

		if err.Error() != "Not allowed remote URL origin: bar.com" {
			t.Fatalf("Invalid error message: %s", err)
		}
	}

	r, _ := http.NewRequest("GET", "http://foo/bar?url=http://bar.com", nil)
	w := httptest.NewRecorder()
	fakeHandler(w, r)
}

func TestHttpImageSourceForwardAuthHeader(t *testing.T) {
	cases := []string{
		"X-Forward-Authorization",
		"Authorization",
	}

	for _, header := range cases {
		r, _ := http.NewRequest("GET", "http://foo/bar?url=http://bar.com", nil)
		r.Header.Set(header, "foobar")

		source := &HttpImageSource{&SourceConfig{AuthForwarding: true}}
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		oreq := &http.Request{Header: make(http.Header)}
		source.setAuthorizationHeader(oreq, r)

		if oreq.Header.Get("Authorization") != "foobar" {
			t.Fatal("Missmatch Authorization header")
		}
	}
}

func TestHttpImageSourceError(t *testing.T) {
	var err error

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Not found"))
	}))
	defer ts.Close()

	source := NewHttpImageSource(&SourceConfig{})
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		_, err = source.GetImage(r)
		if err == nil {
			t.Fatalf("Server response should not be valid: %s", err)
		}
	}

	r, _ := http.NewRequest("GET", "http://foo/bar?url="+ts.URL, nil)
	w := httptest.NewRecorder()
	fakeHandler(w, r)
}

func TestHttpImageSourceExceedsMaximumAllowedLength(t *testing.T) {
	var body []byte
	var err error

	buf, _ := ioutil.ReadFile(fixture1024Bytes)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(buf)
	}))
	defer ts.Close()

	source := NewHttpImageSource(&SourceConfig{
		MaxAllowedSize: 1023,
	})
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		body, err = source.GetImage(r)
		if err == nil {
			t.Fatalf("It should not allow a request to image exceeding maximum allowed size: %s", err)
		}
		w.Write(body)
	}

	r, _ := http.NewRequest("GET", "http://foo/bar?url="+ts.URL, nil)
	w := httptest.NewRecorder()
	fakeHandler(w, r)
}
