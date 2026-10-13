package main

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/h2non/bimg"
)

// defaultSpoolThreshold is the size up to which a request body is read into memory by
// default. A larger one is written to a temporary file and mapped, so that its pages are
// backed by the file and can be reclaimed, instead of growing the heap by a multiple of
// the image.
const defaultSpoolThreshold = 4 << 20

func noRelease() {}

func init() {
	// libvips caches operations by their arguments, and for a loader that is the address
	// of the image buffer. A buffer is unmapped or freed once its request is done, so a
	// later one at the same address could be served a cached result that is not its own.
	bimg.VipsCacheSetMax(0)
	bimg.VipsCacheSetMaxMem(0)
}

// readSpooled reads r whole. Up to threshold bytes it returns an ordinary buffer. Beyond
// it, it returns a read-only mapping of a temporary file in dir, or in the system's
// temporary directory if dir is empty, valid until release is called.
func readSpooled(r io.Reader, dir string, threshold int64) ([]byte, func(), error) {
	head, err := ioutil.ReadAll(io.LimitReader(r, threshold+1))
	if err != nil {
		return nil, noRelease, err
	}
	if int64(len(head)) <= threshold {
		return head, noRelease, nil
	}

	file, err := ioutil.TempFile(dir, "imaginary-")
	if err != nil {
		return nil, noRelease, err
	}
	// Unlinked at once, so that its space is freed when it is closed and unmapped, also
	// when the process dies. A system that cannot unlink an open file removes it after.
	if name := file.Name(); os.Remove(name) != nil {
		defer os.Remove(name)
	}
	defer file.Close()

	if _, err = file.Write(head); err != nil {
		return nil, noRelease, err
	}
	if _, err = io.Copy(file, r); err != nil {
		return nil, noRelease, err
	}

	return mapFile(file)
}
