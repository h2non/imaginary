package main

import (
	"github.com/daddye/vips"
)

func Resize(imageBuf []byte) ([]byte, error) {
	options := vips.Options{
		Width:        562,
		Height:       562,
		Crop:         true,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      100,
	}

	return vips.Resize(imageBuf, options)
}
