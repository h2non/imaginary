package main

import "gopkg.in/h2non/bimg.v0"

// ImageOptions represent all the supported image transformation params as first level members
type ImageOptions struct {
	Width       int
	Height      int
	AreaWidth   int
	AreaHeight  int
	Quality     int
	Compression int
	Rotate      int
	Top         int
	Left        int
	Margin      int
	Factor      int
	DPI         int
	TextWidth   int
	Force       bool
	NoCrop      bool
	NoReplicate bool
	NoRotation  bool
	NoProfile   bool
	Opacity     float32
	Text        string
	Font        string
	Type        string
	Color       []uint8
	Gravity     bimg.Gravity
	Colorspace  bimg.Interpretation
}

// BimgOptions creates a new bimg compatible options struct mapping the fields properly
func BimgOptions(o ImageOptions) bimg.Options {
	debug("Options: %#v", o)
	return bimg.Options{
		Width:          o.Width,
		Height:         o.Height,
		Quality:        o.Quality,
		Compression:    o.Compression,
		NoAutoRotate:   o.NoRotation,
		NoProfile:      o.NoProfile,
		Force:          o.Force,
		Gravity:        o.Gravity,
		Interpretation: o.Colorspace,
		Type:           ImageType(o.Type),
		Rotate:         bimg.Angle(o.Rotate),
	}
}
