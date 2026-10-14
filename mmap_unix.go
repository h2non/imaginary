//go:build !windows
// +build !windows

package main

import (
	"errors"
	"os"
	"syscall"
)

// mapFile maps file read-only. The mapping stays valid after the file is closed, until
// release is called.
func mapFile(file *os.File) ([]byte, func(), error) {
	info, err := file.Stat()
	if err != nil {
		return nil, noRelease, err
	}
	size := info.Size()
	if size == 0 {
		return []byte{}, noRelease, nil
	}
	if int64(int(size)) != size {
		return nil, noRelease, errors.New("image too large to map")
	}

	buf, err := syscall.Mmap(int(file.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, noRelease, err
	}

	return buf, func() { _ = syscall.Munmap(buf) }, nil
}
