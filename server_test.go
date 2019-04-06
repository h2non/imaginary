package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"gopkg.in/h2non/bimg.v1"
)

func TestIndex(t *testing.T) {
	ts := testServer(indexController)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Invalid response status: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(body), "imaginary") == false {
		t.Fatalf("Invalid body response: %s", body)
	}
}

func TestCrop(t *testing.T) {
	ts := testServer(controller(Crop))
	buf := readFile("large.jpg")
	url := ts.URL + "?width=300"
	defer ts.Close()

	res, err := http.Post(url, "image/jpeg", buf)
	if err != nil {
		t.Fatal("Cannot perform the request")
	}

	if res.StatusCode != 200 {
		t.Fatalf("Invalid response status: %s", res.Status)
	}

	if res.Header.Get("Content-Length") == "" {
		t.Fatal("Empty content length response")
	}

	image, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(image) == 0 {
		t.Fatalf("Empty response body")
	}

	err = assertSize(image, 300, 1080)
	if err != nil {
		t.Error(err)
	}

	if bimg.DetermineImageTypeName(image) != "jpeg" {
		t.Fatalf("Invalid image type")
	}
}

func TestResize(t *testing.T) {
	ts := testServer(controller(Resize))
	buf := readFile("large.jpg")
	url := ts.URL + "?width=300&nocrop=false"
	defer ts.Close()

	res, err := http.Post(url, "image/jpeg", buf)
	if err != nil {
		t.Fatal("Cannot perform the request")
	}

	if res.StatusCode != 200 {
		t.Fatalf("Invalid response status: %s", res.Status)
	}

	image, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(image) == 0 {
		t.Fatalf("Empty response body")
	}

	err = assertSize(image, 300, 1080)
	if err != nil {
		t.Error(err)
	}

	if bimg.DetermineImageTypeName(image) != "jpeg" {
		t.Fatalf("Invalid image type")
	}
}

func TestEnlarge(t *testing.T) {
	ts := testServer(controller(Enlarge))
	buf := readFile("large.jpg")
	url := ts.URL + "?width=300&height=200"
	defer ts.Close()

	res, err := http.Post(url, "image/jpeg", buf)
	if err != nil {
		t.Fatal("Cannot perform the request")
	}

	if res.StatusCode != 200 {
		t.Fatalf("Invalid response status: %s", res.Status)
	}

	image, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(image) == 0 {
		t.Fatalf("Empty response body")
	}

	err = assertSize(image, 300, 200)
	if err != nil {
		t.Error(err)
	}

	if bimg.DetermineImageTypeName(image) != "jpeg" {
		t.Fatalf("Invalid image type")
	}
}

func TestExtract(t *testing.T) {
	ts := testServer(controller(Extract))
	buf := readFile("large.jpg")
	url := ts.URL + "?top=100&left=100&areawidth=200&areaheight=120"
	defer ts.Close()

	res, err := http.Post(url, "image/jpeg", buf)
	if err != nil {
		t.Fatal("Cannot perform the request")
	}

	if res.StatusCode != 200 {
		t.Fatalf("Invalid response status: %s", res.Status)
	}

	image, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(image) == 0 {
		t.Fatalf("Empty response body")
	}

	err = assertSize(image, 200, 120)
	if err != nil {
		t.Error(err)
	}

	if bimg.DetermineImageTypeName(image) != "jpeg" {
		t.Fatalf("Invalid image type")
	}
}

func TestTypeAuto(t *testing.T) {
	cases := []struct {
		acceptHeader string
		expected     string
	}{
		{"", "jpeg"},
		{"image/webp,*/*", "webp"},
		{"image/png,*/*", "png"},
		{"image/webp;q=0.8,image/jpeg", "webp"},
		{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8", "webp"}, // Chrome
	}

	for _, test := range cases {
		ts := testServer(controller(Crop))
		buf := readFile("large.jpg")
		url := ts.URL + "?width=300&type=auto"
		defer ts.Close()

		req, _ := http.NewRequest(http.MethodPost, url, buf)
		req.Header.Add("Content-Type", "image/jpeg")
		req.Header.Add("Accept", test.acceptHeader)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("Cannot perform the request")
		}

		if res.StatusCode != 200 {
			t.Fatalf("Invalid response status: %s", res.Status)
		}

		if res.Header.Get("Content-Length") == "" {
			t.Fatal("Empty content length response")
		}

		image, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if len(image) == 0 {
			t.Fatalf("Empty response body")
		}

		err = assertSize(image, 300, 1080)
		if err != nil {
			t.Error(err)
		}

		if bimg.DetermineImageTypeName(image) != test.expected {
			t.Fatalf("Invalid image type")
		}

		if res.Header.Get("Vary") != "Accept" {
			t.Fatal("Vary header not set correctly")
		}
	}
}

func TestFit(t *testing.T) {
	var err error

	buf := readFile("large.jpg")
	original, _ := ioutil.ReadAll(buf)
	err = assertSize(original, 1920, 1080)
	if err != nil {
		t.Errorf("Reference image expecations weren't met")
	}

	ts := testServer(controller(Fit))
	url := ts.URL + "?width=300&height=300"
	defer ts.Close()

	res, err := http.Post(url, "image/jpeg", bytes.NewReader(original))
	if err != nil {
		t.Fatal("Cannot perform the request")
	}

	if res.StatusCode != 200 {
		t.Fatalf("Invalid response status: %s", res.Status)
	}

	image, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(image) == 0 {
		t.Fatalf("Empty response body")
	}

	// The reference image has a ratio of 1.778, this should produce a height of 168.75
	err = assertSize(image, 300, 169)
	if err != nil {
		t.Error(err)
	}

	if bimg.DetermineImageTypeName(image) != "jpeg" {
		t.Fatalf("Invalid image type")
	}
}

func TestRemoteHTTPSource(t *testing.T) {
	opts := ServerOptions{EnableURLSource: true}
	fn := ImageMiddleware(opts)(Crop)
	LoadSources(opts)

	tsImage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		buf, _ := ioutil.ReadFile("testdata/large.jpg")
		w.Write(buf)
	}))
	defer tsImage.Close()

	ts := httptest.NewServer(fn)
	url := ts.URL + "?width=200&height=200&url=" + tsImage.URL
	defer ts.Close()

	res, err := http.Get(url)
	if err != nil {
		t.Fatal("Cannot perform the request")
	}
	if res.StatusCode != 200 {
		t.Fatalf("Invalid response status: %d", res.StatusCode)
	}

	image, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(image) == 0 {
		t.Fatalf("Empty response body")
	}

	err = assertSize(image, 200, 200)
	if err != nil {
		t.Error(err)
	}

	if bimg.DetermineImageTypeName(image) != "jpeg" {
		t.Fatalf("Invalid image type")
	}
}

func TestInvalidRemoteHTTPSource(t *testing.T) {
	opts := ServerOptions{EnableURLSource: true}
	fn := ImageMiddleware(opts)(Crop)
	LoadSources(opts)

	tsImage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(400)
	}))
	defer tsImage.Close()

	ts := httptest.NewServer(fn)
	url := ts.URL + "?width=200&height=200&url=" + tsImage.URL
	defer ts.Close()

	res, err := http.Get(url)
	if err != nil {
		t.Fatal("Request failed")
	}
	if res.StatusCode != 400 {
		t.Fatalf("Invalid response status: %d", res.StatusCode)
	}
}

func TestMountDirectory(t *testing.T) {
	opts := ServerOptions{Mount: "testdata"}
	fn := ImageMiddleware(opts)(Crop)
	LoadSources(opts)

	ts := httptest.NewServer(fn)
	url := ts.URL + "?width=200&height=200&file=large.jpg"
	defer ts.Close()

	res, err := http.Get(url)
	if err != nil {
		t.Fatal("Cannot perform the request")
	}
	if res.StatusCode != 200 {
		t.Fatalf("Invalid response status: %d", res.StatusCode)
	}

	image, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(image) == 0 {
		t.Fatalf("Empty response body")
	}

	err = assertSize(image, 200, 200)
	if err != nil {
		t.Error(err)
	}

	if bimg.DetermineImageTypeName(image) != "jpeg" {
		t.Fatalf("Invalid image type")
	}
}

func TestMountInvalidDirectory(t *testing.T) {
	fn := ImageMiddleware(ServerOptions{Mount: "_invalid_"})(Crop)
	ts := httptest.NewServer(fn)
	url := ts.URL + "?top=100&left=100&areawidth=200&areaheight=120&file=large.jpg"
	defer ts.Close()

	res, err := http.Get(url)
	if err != nil {
		t.Fatal("Cannot perform the request")
	}

	if res.StatusCode != 400 {
		t.Fatalf("Invalid response status: %d", res.StatusCode)
	}
}

func TestMountInvalidPath(t *testing.T) {
	fn := ImageMiddleware(ServerOptions{Mount: "_invalid_"})(Crop)
	ts := httptest.NewServer(fn)
	url := ts.URL + "?top=100&left=100&areawidth=200&areaheight=120&file=../../large.jpg"
	defer ts.Close()

	res, err := http.Get(url)
	if err != nil {
		t.Fatal("Cannot perform the request")
	}

	if res.StatusCode != 400 {
		t.Fatalf("Invalid response status: %s", res.Status)
	}
}

func controller(op Operation) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf, _ := ioutil.ReadAll(r.Body)
		imageHandler(w, r, buf, op, ServerOptions{})
	}
}

func testServer(fn func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(fn))
}

func readFile(file string) io.Reader {
	buf, _ := os.Open(path.Join("testdata", file))
	return buf
}

func assertSize(buf []byte, width, height int) error {
	size, err := bimg.NewImage(buf).Size()
	if err != nil {
		return err
	}
	if size.Width != width || size.Height != height {
		return fmt.Errorf("Invalid image size: %dx%d, expected: %dx%d", size.Width, size.Height, width, height)
	}
	return nil
}
