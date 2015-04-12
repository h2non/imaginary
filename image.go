package main

import (
	"gopkg.in/h2non/bimg.v0"
)

type ImageOptions struct {
	Width       int
	Height      int
	Quality     int
	Compression int
	Rotate      int
	Top         int
	Left        int
	Margin      int
	Factor      int
	Opacity     float64
	Text        string
	Font        string
	Color       string
	Type        string
}

type Operation func([]byte, ImageOptions) ([]byte, error)

func (o Operation) Run(buf []byte, opts ImageOptions) ([]byte, error) {
	return o(buf, opts)
}

func bimgOptions(o ImageOptions) bimg.Options {
	return bimg.Options{
		Width:       o.Width,
		Height:      o.Height,
		Quality:     o.Quality,
		Compression: o.Quality,
		Type:        ImageType(o.Type),
	}
}

func Resize(buf []byte, o ImageOptions) ([]byte, error) {
	if o.Width == 0 || o.Height == 0 {
		return nil, NewError("Missing required params: height, width", BAD_REQUEST)
	}

	opts := bimgOptions(o)
	return Process(buf, opts)
}

func Enlarge(buf []byte, o ImageOptions) ([]byte, error) {
	if o.Width == 0 || o.Height == 0 {
		return nil, NewError("Missing required params: height, width", BAD_REQUEST)
	}

	opts := bimgOptions(o)
	opts.Enlarge = true
	return Process(buf, opts)
}

func Crop(buf []byte, o ImageOptions) ([]byte, error) {
	opts := bimgOptions(o)
	opts.Crop = true
	return Process(buf, opts)
}

func Rotate(buf []byte, o ImageOptions) ([]byte, error) {
	if o.Rotate == 0 {
		return nil, NewError("Missing rotate param", BAD_REQUEST)
	}

	opts := bimgOptions(o)
	opts.Rotate = bimg.Angle(o.Rotate)
	return Process(buf, opts)
}

func Flip(buf []byte, o ImageOptions) ([]byte, error) {
	opts := bimgOptions(o)
	opts.Flip = true
	return Process(buf, opts)
}

func Flop(buf []byte, o ImageOptions) ([]byte, error) {
	opts := bimgOptions(o)
	opts.Flop = true
	return Process(buf, opts)
}

func Thumbnail(buf []byte, o ImageOptions) ([]byte, error) {
	if o.Width == 0 && o.Height == 0 {
		return nil, NewError("Missing required params: width or height", BAD_REQUEST)
	}

	opts := bimgOptions(o)
	return Process(buf, opts)
}

func Zoom(buf []byte, o ImageOptions) ([]byte, error) {
	debug("Options: ")
	if o.Width == 0 || o.Height == 0 || o.Factor == 0 {
		return nil, NewError("Missing required params: width, height, factor", BAD_REQUEST)
	}

	opts := bimgOptions(o)
	//opts.Crop = true
	opts.Zoom = o.Factor
	return Process(buf, opts)
}

func Convert(buf []byte, o ImageOptions) ([]byte, error) {
	if o.Type == "" {
		return nil, NewError("Missing required params: type", BAD_REQUEST)
	}

	opts := bimgOptions(o)
	return Process(buf, opts)
}

func Watermark(buf []byte, o ImageOptions) ([]byte, error) {
	if o.Text == "" {
		return nil, NewError("Missing required params: text", BAD_REQUEST)
	}

	opts := bimgOptions(o)
	opts.Watermark.Text = o.Text
	opts.Watermark.Font = o.Font
	opts.Watermark.Margin = o.Margin
	opts.Watermark.Opacity = float32(o.Opacity)
	return Process(buf, opts)
}

func Extract(buf []byte, o ImageOptions) ([]byte, error) {
	if o.Top == 0 || o.Left == 0 {
		return nil, NewError("Missing required params: top, left", BAD_REQUEST)
	}

	opts := bimgOptions(o)
	return Process(buf, opts)
}

type ImageInfo struct {
	Size        int    `json:"size"`
	Width       int    `json:"width"`
	Format      string `json:"format"`
	Height      int    `json:"height"`
	Orientation int    `json:"orientation"`
	Alpha       bool   `json:"alpha"`
}

func Info(buf []byte, o ImageOptions) ([]byte, error) {
	return []byte{}, nil
}

func Process(buf []byte, opts bimg.Options) ([]byte, error) {
	return bimg.Resize(buf, opts)
}
