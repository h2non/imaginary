package main

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func randomBytes(n int) []byte {
	buf := make([]byte, n)
	rand.New(rand.NewSource(1)).Read(buf)
	return buf
}

func TestReadSpooledSmall(t *testing.T) {
	dir := t.TempDir()
	data := randomBytes(1024)

	buf, release, err := readSpooled(bytes.NewReader(data), dir, defaultSpoolThreshold)
	if err != nil {
		t.Fatal(err)
	}
	defer release()

	if !bytes.Equal(buf, data) {
		t.Error("small body read wrongly")
	}
	if files, _ := ioutil.ReadDir(dir); len(files) != 0 {
		t.Errorf("small body spooled: %d files", len(files))
	}
}

func TestReadSpooledLarge(t *testing.T) {
	dir := t.TempDir()
	data := randomBytes(defaultSpoolThreshold + 1000)

	buf, release, err := readSpooled(bytes.NewReader(data), dir, defaultSpoolThreshold)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf, data) {
		t.Error("large body read wrongly")
	}
	// The temporary file is unlinked while it is still mapped
	if files, _ := ioutil.ReadDir(dir); len(files) != 0 {
		t.Errorf("temporary file left in %s: %d files", dir, len(files))
	}
	release()
}

func TestReadSpooledAtThreshold(t *testing.T) {
	data := randomBytes(defaultSpoolThreshold)

	buf, release, err := readSpooled(bytes.NewReader(data), t.TempDir(), defaultSpoolThreshold)
	if err != nil {
		t.Fatal(err)
	}
	defer release()

	if !bytes.Equal(buf, data) {
		t.Error("body at the threshold read wrongly")
	}
}

func TestBodyImageSourceLarge(t *testing.T) {
	data := randomBytes(defaultSpoolThreshold * 2)
	source := NewBodyImageSource(&SourceConfig{TmpDir: t.TempDir(), SpoolThreshold: defaultSpoolThreshold}).(ReleasingImageSource)

	r, _ := http.NewRequest(http.MethodPost, "http://foo/bar", bytes.NewReader(data))
	buf, release, err := source.GetImageReleasing(r)
	if err != nil {
		t.Fatal(err)
	}
	defer release()

	if !bytes.Equal(buf, data) {
		t.Error("large body read wrongly")
	}
}

func TestBodyImageSourceGetImageOwnsCopy(t *testing.T) {
	data := randomBytes(defaultSpoolThreshold * 2)
	source := NewBodyImageSource(&SourceConfig{TmpDir: t.TempDir(), SpoolThreshold: defaultSpoolThreshold})

	r, _ := http.NewRequest(http.MethodPost, "http://foo/bar", bytes.NewReader(data))
	buf, err := source.GetImage(r)
	if err != nil {
		t.Fatal(err)
	}

	// Still readable, as the mapping is released before GetImage returns
	if !bytes.Equal(buf, data) {
		t.Error("large body read wrongly")
	}
}

func TestReadSpooledThreshold(t *testing.T) {
	// A directory that does not exist makes spooling fail, which shows which path was taken
	missing := filepath.Join(t.TempDir(), "missing")
	data := randomBytes(2048)

	cases := []struct {
		threshold int64
		size      int
		spooled   bool
	}{
		{1024, 1024, false},
		{1024, 1025, true},
		{4096, 2048, false},
		{0, 1, true},
	}
	for _, c := range cases {
		buf, release, err := readSpooled(bytes.NewReader(data[:c.size]), missing, c.threshold)
		release()
		if c.spooled && err == nil {
			t.Errorf("threshold %d, size %d: read into memory, want spooled", c.threshold, c.size)
		}
		if !c.spooled && (err != nil || !bytes.Equal(buf, data[:c.size])) {
			t.Errorf("threshold %d, size %d: %v, want read into memory", c.threshold, c.size, err)
		}
	}
}

func TestMapFile(t *testing.T) {
	data := randomBytes(defaultSpoolThreshold + 1)
	file, err := ioutil.TempFile(t.TempDir(), "map-")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = file.Write(data); err != nil {
		t.Fatal(err)
	}

	buf, release, err := mapFile(file)
	file.Close()
	if err != nil {
		t.Fatal(err)
	}
	defer release()

	// The mapping outlives the file descriptor
	if !bytes.Equal(buf, data) {
		t.Error("mapped file read wrongly")
	}
}

func TestImageControllerLargeBody(t *testing.T) {
	buf, err := ioutil.ReadFile(fixtureFile)
	if err != nil {
		t.Fatal(err)
	}
	// Padded past the threshold, which the JPEG decoder ignores after the end marker
	body := append(buf, make([]byte, defaultSpoolThreshold)...)

	opts := ServerOptions{MaxAllowedPixels: 18.0, TmpDir: t.TempDir(), SpoolThreshold: defaultSpoolThreshold}
	LoadSources(opts)
	ts := httptest.NewServer(ImageMiddleware(opts)(Crop))
	defer ts.Close()

	res, err := http.Post(ts.URL+"?width=200&height=200", "image/jpeg", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	image, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status %d: %s", res.StatusCode, image)
	}
	if err = assertSize(image, 200, 200); err != nil {
		t.Error(err)
	}
}
