package main

import (
	"io/ioutil"
	"testing"
)

func TestImageResize(t *testing.T) {
	t.Run("Width and Height defined", func(t *testing.T) {
		opts := ImageOptions{Width: 300, Height: 300}
		buf, _ := ioutil.ReadAll(readFile("imaginary.jpg"))

		img, err := Resize(buf, opts)
		if err != nil {
			t.Errorf("Cannot process image: %s", err)
		}
		if img.Mime != "image/jpeg" {
			t.Error("Invalid image MIME type")
		}
		if assertSize(img.Body, opts.Width, opts.Height) != nil {
			t.Errorf("Invalid image size, expected: %dx%d", opts.Width, opts.Height)
		}
	})

	t.Run("Width defined", func(t *testing.T) {
		opts := ImageOptions{Width: 300}
		buf, _ := ioutil.ReadAll(readFile("imaginary.jpg"))

		img, err := Resize(buf, opts)
		if err != nil {
			t.Errorf("Cannot process image: %s", err)
		}
		if img.Mime != "image/jpeg" {
			t.Error("Invalid image MIME type")
		}
		if err := assertSize(img.Body, 300, 404); err != nil {
			t.Error(err)
		}
	})

	t.Run("Width defined with NoCrop=false", func(t *testing.T) {
		opts := ImageOptions{Width: 300, NoCrop: false, IsDefinedField: IsDefinedField{NoCrop: true}}
		buf, _ := ioutil.ReadAll(readFile("imaginary.jpg"))

		img, err := Resize(buf, opts)
		if err != nil {
			t.Errorf("Cannot process image: %s", err)
		}
		if img.Mime != "image/jpeg" {
			t.Error("Invalid image MIME type")
		}

		// The original image is 550x740
		if err := assertSize(img.Body, 300, 740); err != nil {
			t.Error(err)
		}
	})

	t.Run("Width defined with NoCrop=true", func(t *testing.T) {
		opts := ImageOptions{Width: 300, NoCrop: true, IsDefinedField: IsDefinedField{NoCrop: true}}
		buf, _ := ioutil.ReadAll(readFile("imaginary.jpg"))

		img, err := Resize(buf, opts)
		if err != nil {
			t.Errorf("Cannot process image: %s", err)
		}
		if img.Mime != "image/jpeg" {
			t.Error("Invalid image MIME type")
		}

		// The original image is 550x740
		if err := assertSize(img.Body, 300, 404); err != nil {
			t.Error(err)
		}
	})

}

func TestImageFit(t *testing.T) {
	opts := ImageOptions{Width: 300, Height: 300}
	buf, _ := ioutil.ReadAll(readFile("imaginary.jpg"))

	img, err := Fit(buf, opts)
	if err != nil {
		t.Errorf("Cannot process image: %s", err)
	}
	if img.Mime != "image/jpeg" {
		t.Error("Invalid image MIME type")
	}
	// 550x740 -> 222.9x300
	if assertSize(img.Body, 223, 300) != nil {
		t.Errorf("Invalid image size, expected: %dx%d", opts.Width, opts.Height)
	}
}

func TestImagePipelineOperations(t *testing.T) {
	width, height := 300, 260

	operations := PipelineOperations{
		PipelineOperation{
			Name: "crop",
			Params: map[string]interface{}{
				"width":  width,
				"height": height,
			},
		},
		PipelineOperation{
			Name: "convert",
			Params: map[string]interface{}{
				"type": "webp",
			},
		},
	}

	opts := ImageOptions{Operations: operations}
	buf, _ := ioutil.ReadAll(readFile("imaginary.jpg"))

	img, err := Pipeline(buf, opts)
	if err != nil {
		t.Errorf("Cannot process image: %s", err)
	}
	if img.Mime != "image/webp" {
		t.Error("Invalid image MIME type")
	}
	if assertSize(img.Body, width, height) != nil {
		t.Errorf("Invalid image size, expected: %dx%d", width, height)
	}
}

func TestCalculateDestinationFitDimension(t *testing.T) {
	cases := []struct {
		// Image
		imageWidth  int
		imageHeight int

		// User parameter
		optionWidth  int
		optionHeight int

		// Expect
		fitWidth  int
		fitHeight int
	}{

		// Leading Width
		{1280, 1000, 710, 9999, 710, 555},
		{1279, 1000, 710, 9999, 710, 555},
		{900, 500, 312, 312, 312, 173}, // rounding down
		{900, 500, 313, 313, 313, 174}, // rounding up

		// Leading height
		{1299, 2000, 710, 999, 649, 999},
		{1500, 2000, 710, 999, 710, 947},
	}

	for _, tc := range cases {
		fitWidth, fitHeight := calculateDestinationFitDimension(tc.imageWidth, tc.imageHeight, tc.optionWidth, tc.optionHeight)
		if fitWidth != tc.fitWidth || fitHeight != tc.fitHeight {
			t.Errorf(
				"Fit dimensions calculation failure\nExpected : %d/%d (width/height)\nActual   : %d/%d (width/height)\n%+v",
				tc.fitWidth, tc.fitHeight, fitWidth, fitHeight, tc,
			)
		}
	}

}
