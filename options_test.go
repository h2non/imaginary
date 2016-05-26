package main

import "testing"

func TestBimgOptions(t *testing.T) {
	imgOpts := ImageOptions{
		Width:  500,
		Height: 600,
	}
	opts := BimgOptions(imgOpts)

	if opts.Width != imgOpts.Width || opts.Height != imgOpts.Height {
		t.Error("Invalid width and height")
	}
}
