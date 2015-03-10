package main

import (
	"github.com/daddye/vips"
)

type Options struct {
	*vips.Options
	Width  int
	Height int
	Crop   bool
}
