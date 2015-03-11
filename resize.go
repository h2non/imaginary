package main

import (
	"github.com/daddye/vips"
)

func Resize(imageBuf []byte) ([]byte, error) {
	options := vips.Options{
		Width:        300,
		Height:       240,
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      95,
	}

	return vips.Resize(imageBuf, options)
}
