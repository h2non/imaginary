package main

import (
	"errors"
	"net/http"
	"strconv"
)

type ImageOptions struct {
	Width   int
	Height  int
	Quality int
}

type Image struct {
	Body []byte
}

func validateParams(r *http.Request) (*ImageOptions, error) {
	query := r.URL.Query()
	width, _ := strconv.Atoi(query.Get("width"))
	height, _ := strconv.Atoi(query.Get("height"))
	quality, _ := strconv.Atoi(query.Get("quality"))

	if width == 0 || height == 0 {
		return nil, errors.New("Missing required height and width params")
	}

	return &ImageOptions{
		Width:   width,
		Height:  height,
		Quality: quality,
	}, nil
}
