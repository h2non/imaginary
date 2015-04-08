package main

import (
	"gopkg.in/h2non/bimg.v0"
)

func Resize(imageBuf []byte) ([]byte, error) {
	options := bimg.Options{
		Width:   562,
		Height:  562,
		Crop:    true,
		Quality: 95,
	}

	return bimg.Resize(imageBuf, options)
}
