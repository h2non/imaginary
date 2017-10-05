package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/h2non/bimg.v1"
)

// OperationsMap defines the allowed image transformation operations listed by name.
// Used for pipeline image processing.
var OperationsMap = map[string]Operation{
	"crop":      Crop,
	"resize":    Resize,
	"enlarge":   Enlarge,
	"extract":   Extract,
	"rotate":    Rotate,
	"flip":      Flip,
	"flop":      Flop,
	"thumbnail": Thumbnail,
	"zoom":      Zoom,
	"convert":   Convert,
	"watermark": Watermark,
	"blur":      GaussianBlur,
	"smartcrop": SmartCrop,
}

// Image stores an image binary buffer and its MIME type
type Image struct {
	Body []byte
	Mime string
}

// Operation implements an image transformation runnable interface
type Operation func([]byte, ImageOptions) (Image, error)

// Run performs the image transformation
func (o Operation) Run(buf []byte, opts ImageOptions) (Image, error) {
	return o(buf, opts)
}

// ImageInfo represents an image details and additional metadata
type ImageInfo struct {
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Type        string `json:"type"`
	Space       string `json:"space"`
	Alpha       bool   `json:"hasAlpha"`
	Profile     bool   `json:"hasProfile"`
	Channels    int    `json:"channels"`
	Orientation int    `json:"orientation"`
}

func Info(buf []byte, o ImageOptions) (Image, error) {
	// We're not handling an image here, but we reused the struct.
	// An interface will be definitively better here.
	image := Image{Mime: "application/json"}

	meta, err := bimg.Metadata(buf)
	if err != nil {
		return image, NewError("Cannot retrieve image metadata: %s"+err.Error(), BadRequest)
	}

	info := ImageInfo{
		Width:       meta.Size.Width,
		Height:      meta.Size.Height,
		Type:        meta.Type,
		Space:       meta.Space,
		Alpha:       meta.Alpha,
		Profile:     meta.Profile,
		Channels:    meta.Channels,
		Orientation: meta.Orientation,
	}

	body, _ := json.Marshal(info)
	image.Body = body

	return image, nil
}

func Resize(buf []byte, o ImageOptions) (Image, error) {
	if o.Width == 0 && o.Height == 0 {
		return Image{}, NewError("Missing required param: height or width", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Embed = true

	if o.NoCrop == false {
		opts.Crop = true
	}

	return Process(buf, opts)
}

func Enlarge(buf []byte, o ImageOptions) (Image, error) {
	if o.Width == 0 || o.Height == 0 {
		return Image{}, NewError("Missing required params: height, width", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Enlarge = true

	if o.NoCrop == false {
		opts.Crop = true
	}

	return Process(buf, opts)
}

func Extract(buf []byte, o ImageOptions) (Image, error) {
	if o.AreaWidth == 0 || o.AreaHeight == 0 {
		return Image{}, NewError("Missing required params: areawidth or areaheight", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Top = o.Top
	opts.Left = o.Left
	opts.AreaWidth = o.AreaWidth
	opts.AreaHeight = o.AreaHeight

	return Process(buf, opts)
}

func Crop(buf []byte, o ImageOptions) (Image, error) {
	if o.Width == 0 && o.Height == 0 {
		return Image{}, NewError("Missing required param: height or width", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Crop = true
	return Process(buf, opts)
}

func SmartCrop(buf []byte, o ImageOptions) (Image, error) {
	if o.Width == 0 && o.Height == 0 {
		return Image{}, NewError("Missing required param: height or width", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Crop = true
	opts.Gravity = bimg.GravitySmart
	return Process(buf, opts)
}

func Rotate(buf []byte, o ImageOptions) (Image, error) {
	if o.Rotate == 0 {
		return Image{}, NewError("Missing required param: rotate", BadRequest)
	}

	opts := BimgOptions(o)
	return Process(buf, opts)
}

func Flip(buf []byte, o ImageOptions) (Image, error) {
	opts := BimgOptions(o)
	opts.Flip = true
	return Process(buf, opts)
}

func Flop(buf []byte, o ImageOptions) (Image, error) {
	opts := BimgOptions(o)
	opts.Flop = true
	return Process(buf, opts)
}

func Thumbnail(buf []byte, o ImageOptions) (Image, error) {
	if o.Width == 0 && o.Height == 0 {
		return Image{}, NewError("Missing required params: width or height", BadRequest)
	}

	return Process(buf, BimgOptions(o))
}

func Zoom(buf []byte, o ImageOptions) (Image, error) {
	if o.Factor == 0 {
		return Image{}, NewError("Missing required param: factor", BadRequest)
	}

	opts := BimgOptions(o)

	if o.Top > 0 || o.Left > 0 {
		if o.AreaWidth == 0 && o.AreaHeight == 0 {
			return Image{}, NewError("Missing required params: areawidth, areaheight", BadRequest)
		}

		opts.Top = o.Top
		opts.Left = o.Left
		opts.AreaWidth = o.AreaWidth
		opts.AreaHeight = o.AreaHeight

		if o.NoCrop == false {
			opts.Crop = true
		}
	}

	opts.Zoom = o.Factor
	return Process(buf, opts)
}

func Convert(buf []byte, o ImageOptions) (Image, error) {
	if o.Type == "" {
		return Image{}, NewError("Missing required param: type", BadRequest)
	}
	if ImageType(o.Type) == bimg.UNKNOWN {
		return Image{}, NewError("Invalid image type: "+o.Type, BadRequest)
	}
	opts := BimgOptions(o)

	return Process(buf, opts)
}

func Watermark(buf []byte, o ImageOptions) (Image, error) {
	if o.Text == "" {
		return Image{}, NewError("Missing required param: text", BadRequest)
	}

	opts := BimgOptions(o)
	opts.Watermark.DPI = o.DPI
	opts.Watermark.Text = o.Text
	opts.Watermark.Font = o.Font
	opts.Watermark.Margin = o.Margin
	opts.Watermark.Width = o.TextWidth
	opts.Watermark.Opacity = o.Opacity
	opts.Watermark.NoReplicate = o.NoReplicate

	if len(o.Color) > 2 {
		opts.Watermark.Background = bimg.Color{o.Color[0], o.Color[1], o.Color[2]}
	}

	return Process(buf, opts)
}

func GaussianBlur(buf []byte, o ImageOptions) (Image, error) {
	if o.Sigma == 0 && o.MinAmpl == 0 {
		return Image{}, NewError("Missing required param: sigma or minampl", BadRequest)
	}
	opts := BimgOptions(o)
	return Process(buf, opts)
}

func Pipeline(buf []byte, o ImageOptions) (Image, error) {
	if len(o.Operations) == 0 {
		return Image{}, NewError("Missing or invalid pipeline operations JSON", BadRequest)
	}
	if len(o.Operations) > 10 {
		return Image{}, NewError("Maximum allowed pipeline operations exceeded", BadRequest)
	}

	// Validate and built operations
	for i, operation := range o.Operations {
		// Normalize operation name
		name := strings.TrimSpace(strings.ToLower(operation.Name))

		// Validate supported operation name
		var exists bool
		if operation.Operation, exists = OperationsMap[operation.Name]; !exists {
			return Image{}, NewError(fmt.Sprintf("Unsupported operation name: %s", name), BadRequest)
		}

		// Parse and construct operation options
		operation.ImageOptions = readMapParams(operation.Params)

		// Mutate list by value
		o.Operations[i] = operation
	}

	var image Image
	var err error

	// Reduce image by running multiple operations
	image = Image{Body: buf}
	for _, operation := range o.Operations {
		var curImage Image
		curImage, err = operation.Operation(image.Body, operation.ImageOptions)
		if err != nil && !operation.IgnoreFailure {
			return Image{}, err
		}
		if operation.IgnoreFailure {
			err = nil
		}
		if err == nil {
			image = curImage
		}
	}

	return image, err
}

func Process(buf []byte, opts bimg.Options) (out Image, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch value := r.(type) {
			case error:
				err = value
			case string:
				err = errors.New(value)
			default:
				err = errors.New("libvips internal error")
			}
			out = Image{}
		}
	}()

	buf, err = bimg.Resize(buf, opts)
	if err != nil {
		return Image{}, err
	}

	mime := GetImageMimeType(bimg.DetermineImageType(buf))
	return Image{Body: buf, Mime: mime}, nil
}
