package main

import "gopkg.in/h2non/bimg.v1"

// ImageOptions represent all the supported image transformation params as first level members
type ImageOptions struct {
	IsDefinedField

	Width         int
	Height        int
	AreaWidth     int
	AreaHeight    int
	Quality       int
	Compression   int
	Rotate        int
	Top           int
	Left          int
	Margin        int
	Factor        int
	DPI           int
	TextWidth     int
	Flip          bool
	Flop          bool
	Force         bool
	Embed         bool
	NoCrop        bool
	NoReplicate   bool
	NoRotation    bool
	NoProfile     bool
	StripMetadata bool
	Opacity       float32
	Sigma         float64
	MinAmpl       float64
	Text          string
	Image         string
	Font          string
	Type          string
	Color         []uint8
	Background    []uint8
	Extend        bimg.Extend
	Gravity       bimg.Gravity
	Colorspace    bimg.Interpretation
	Operations    PipelineOperations
}

// IsDefinedField holds boolean ImageOptions fields. If true it means the field was specified in the request. This
// metadata allows for sane usage of default (false) values.
type IsDefinedField struct {
	Flip          bool
	Flop          bool
	Force         bool
	Embed         bool
	NoCrop        bool
	NoReplicate   bool
	NoRotation    bool
	NoProfile     bool
	StripMetadata bool
}

// PipelineOperation represents the structure for an operation field.
type PipelineOperation struct {
	Name          string                 `json:"operation"`
	IgnoreFailure bool                   `json:"ignore_failure"`
	Params        map[string]interface{} `json:"params"`
	ImageOptions  ImageOptions           `json:"-"`
	Operation     Operation              `json:"-"`
}

// PipelineOperations defines the expected interface for a list of operations.
type PipelineOperations []PipelineOperation

// BimgOptions creates a new bimg compatible options struct mapping the fields properly
func BimgOptions(o ImageOptions) bimg.Options {
	opts := bimg.Options{
		Width:          o.Width,
		Height:         o.Height,
		Flip:           o.Flip,
		Flop:           o.Flop,
		Quality:        o.Quality,
		Compression:    o.Compression,
		NoAutoRotate:   o.NoRotation,
		NoProfile:      o.NoProfile,
		Force:          o.Force,
		Gravity:        o.Gravity,
		Embed:          o.Embed,
		Extend:         o.Extend,
		Interpretation: o.Colorspace,
		StripMetadata:  o.StripMetadata,
		Type:           ImageType(o.Type),
		Rotate:         bimg.Angle(o.Rotate),
	}

	if len(o.Background) != 0 {
		opts.Background = bimg.Color{R: o.Background[0], G: o.Background[1], B: o.Background[2]}
	}

	if o.Sigma > 0 || o.MinAmpl > 0 {
		opts.GaussianBlur = bimg.GaussianBlur{
			Sigma:   o.Sigma,
			MinAmpl: o.MinAmpl,
		}
	}

	return opts
}
