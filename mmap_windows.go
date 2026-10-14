package main

import (
	"io"
	"io/ioutil"
	"os"
)

// mapFile reads file into memory, as there is no mapping here.
func mapFile(file *os.File) ([]byte, func(), error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, noRelease, err
	}
	buf, err := ioutil.ReadAll(file)
	return buf, noRelease, err
}
