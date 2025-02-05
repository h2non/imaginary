package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestFileSystemImageSource(t *testing.T) {
    var body []byte
    var err error
    const fixtureFile = "testdata/large.jpg"

    source := NewFileSystemImageSource(&SourceConfig{MountPath: "testdata"})
    fakeHandler := func(w http.ResponseWriter, r *http.Request) {
        if !source.Matches(r) {
            t.Fatal("Cannot match the request")
        }

        // Log the file path being accessed
        fmt.Printf("Request URL: %s\n", r.URL.String())

        body, err = source.GetImage(r)
        if err != nil {
            t.Fatalf("Error while reading the body: %s", err)
        }
        _, _ = w.Write(body)
    }

    r, _ := http.NewRequest(http.MethodGet, "http://foo/bar?file=large.jpg", nil)
    w := httptest.NewRecorder()
    fakeHandler(w, r)

    buf, _ := ioutil.ReadFile(fixtureFile)
    if len(body) != len(buf) {
        t.Error("Invalid response body")
    }
}