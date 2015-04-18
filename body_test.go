package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const fixtureFile = "fixtures/large.jpg"

func TestReadBody(t *testing.T) {
	var body []byte
	var err error

	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		body, err = readBody(r)
		if err != nil {
			t.Fatalf("Error while reading the body: %s", err)
		}
		w.Write(body)
	}

	file, _ := os.Open(fixtureFile)
	r, _ := http.NewRequest("POST", "http://foo/bar", file)
	w := httptest.NewRecorder()
	fakeHandler(w, r)

	buf, _ := ioutil.ReadFile(fixtureFile)
	if len(body) != len(buf) {
		t.Error("Invalid response body")
	}
}
