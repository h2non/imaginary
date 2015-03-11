package main

import (
	"github.com/daddye/vips"
)

func NewOptions(width int, height int, quality int) *vips.Options {
	return &vips.Options{
		Width:        width,
		Height:       height,
		Quality:      quality,
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
	}
}
