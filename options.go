package main

import (
	"gopkg.in/h2non/bimg.v0"
)

func NewOptions(width int, height int, quality int) bimg.Options {
	return bimg.Options{
		Width:   width,
		Height:  height,
		Quality: quality,
		Crop:    false,
	}
}
