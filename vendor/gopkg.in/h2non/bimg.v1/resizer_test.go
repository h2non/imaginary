package bimg

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestResize(t *testing.T) {
	options := Options{Width: 800, Height: 600}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != JPEG {
		t.Fatal("Image is not jpeg")
	}

	size, _ := Size(newImg)
	if size.Height != options.Height || size.Width != options.Width {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	Write("testdata/test_out.jpg", newImg)
}

func TestResizeVerticalImage(t *testing.T) {
	tests := []Options{
		{Width: 800, Height: 600},
		{Width: 1000, Height: 1000},
		{Width: 1000, Height: 1500},
		{Width: 1000},
		{Height: 1500},
		{Width: 100, Height: 50},
		{Width: 2000, Height: 2000},
		{Width: 500, Height: 1000},
		{Width: 500},
		{Height: 500},
		{Crop: true, Width: 500, Height: 1000},
		{Crop: true, Enlarge: true, Width: 2000, Height: 1400},
		{Enlarge: true, Force: true, Width: 2000, Height: 2000},
		{Force: true, Width: 2000, Height: 2000},
	}

	bufJpeg, err := Read("testdata/vertical.jpg")
	if err != nil {
		t.Fatal(err)
	}
	bufWebp, err := Read("testdata/vertical.webp")
	if err != nil {
		t.Fatal(err)
	}

	images := []struct {
		format ImageType
		buf    []byte
	}{
		{JPEG, bufJpeg},
		{WEBP, bufWebp},
	}

	for _, source := range images {
		for _, options := range tests {
			image, err := Resize(source.buf, options)
			if err != nil {
				t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
			}

			format := DetermineImageType(image)
			if format != source.format {
				t.Fatalf("Image format is invalid. Expected: %#v got %v", ImageTypeName(source.format), ImageTypeName(format))
			}

			size, _ := Size(image)
			if options.Height > 0 && size.Height != options.Height {
				t.Fatalf("Invalid height: %d", size.Height)
			}
			if options.Width > 0 && size.Width != options.Width {
				t.Fatalf("Invalid width: %d", size.Width)
			}

			Write(
				fmt.Sprintf(
					"testdata/test_vertical_%dx%d_out.%s",
					options.Width,
					options.Height,
					ImageTypeName(source.format)),
				image)
		}
	}
}

func TestResizeCustomSizes(t *testing.T) {
	tests := []Options{
		{Width: 800, Height: 600},
		{Width: 1000, Height: 1000},
		{Width: 100, Height: 50},
		{Width: 2000, Height: 2000},
		{Width: 500, Height: 1000},
		{Width: 500},
		{Height: 500},
		{Crop: true, Width: 500, Height: 1000},
		{Crop: true, Enlarge: true, Width: 2000, Height: 1400},
		{Enlarge: true, Force: true, Width: 2000, Height: 2000},
		{Force: true, Width: 2000, Height: 2000},
	}

	bufJpeg, err := Read("testdata/test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	bufWebp, err := Read("testdata/test.webp")
	if err != nil {
		t.Fatal(err)
	}

	images := []struct {
		format ImageType
		buf    []byte
	}{
		{JPEG, bufJpeg},
		{WEBP, bufWebp},
	}

	for _, source := range images {
		for _, options := range tests {
			image, err := Resize(source.buf, options)
			if err != nil {
				t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
			}

			if DetermineImageType(image) != source.format {
				t.Fatalf("Image format is invalid. Expected: %#v", source.format)
			}

			size, _ := Size(image)

			invalidHeight := options.Height > 0 && size.Height != options.Height
			if !options.Crop && invalidHeight {
				t.Fatalf("Invalid height: %d, expected %d", size.Height, options.Height)
			}

			invalidWidth := options.Width > 0 && size.Width != options.Width
			if !options.Crop && invalidWidth {
				t.Fatalf("Invalid width: %d, expected %d", size.Width, options.Width)
			}

			if options.Crop && invalidHeight && invalidWidth {
				t.Fatalf("Invalid width or height: %dx%d, expected %dx%d (crop)", size.Width, size.Height, options.Width, options.Height)
			}
		}
	}
}

func TestResizePrecision(t *testing.T) {
	// see https://github.com/h2non/bimg/issues/99
	img := image.NewGray16(image.Rect(0, 0, 1920, 1080))
	input := &bytes.Buffer{}
	jpeg.Encode(input, img, nil)

	opts := Options{Width: 300}
	newImg, err := Resize(input.Bytes(), opts)
	if err != nil {
		t.Fatalf("Resize(imgData, %#v) error: %#v", opts, err)
	}

	size, _ := Size(newImg)
	if size.Width != opts.Width {
		t.Fatalf("Invalid width: %d", size.Width)
	}
}

func TestRotate(t *testing.T) {
	options := Options{Width: 800, Height: 600, Rotate: 270, Crop: true}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != JPEG {
		t.Error("Image is not jpeg")
	}

	size, _ := Size(newImg)
	if size.Width != options.Width || size.Height != options.Height {
		t.Errorf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	Write("testdata/test_rotate_out.jpg", newImg)
}

func TestInvalidRotateDegrees(t *testing.T) {
	options := Options{Width: 800, Height: 600, Rotate: 111, Crop: true}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != JPEG {
		t.Errorf("Image is not jpeg")
	}

	size, _ := Size(newImg)
	if size.Width != options.Width || size.Height != options.Height {
		t.Errorf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	Write("testdata/test_rotate_invalid_out.jpg", newImg)
}

func TestCorruptedImage(t *testing.T) {
	options := Options{Width: 800, Height: 600}
	buf, _ := Read("testdata/corrupt.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != JPEG {
		t.Fatal("Image is not jpeg")
	}

	size, _ := Size(newImg)
	if size.Height != options.Height || size.Width != options.Width {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	Write("testdata/test_corrupt_out.jpg", newImg)
}

func TestNoColorProfile(t *testing.T) {
	options := Options{Width: 800, Height: 600, NoProfile: true}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	metadata, err := Metadata(newImg)
	if metadata.Profile == true {
		t.Fatal("Invalid profile data")
	}

	size, _ := Size(newImg)
	if size.Height != options.Height || size.Width != options.Width {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}
}

func TestEmbedExtendColor(t *testing.T) {
	options := Options{Width: 400, Height: 600, Crop: false, Embed: true, Extend: ExtendWhite, Background: Color{255, 20, 10}}
	buf, _ := Read("testdata/test_issue.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	size, _ := Size(newImg)
	if size.Height != options.Height || size.Width != options.Width {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	Write("testdata/test_extend_white_out.jpg", newImg)
}

func TestEmbedExtendWithCustomColor(t *testing.T) {
	options := Options{Width: 400, Height: 600, Crop: false, Embed: true, Extend: 5, Background: Color{255, 20, 10}}
	buf, _ := Read("testdata/test_issue.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	size, _ := Size(newImg)
	if size.Height != options.Height || size.Width != options.Width {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	Write("testdata/test_extend_background_out.jpg", newImg)
}

func TestGaussianBlur(t *testing.T) {
	options := Options{Width: 800, Height: 600, GaussianBlur: GaussianBlur{Sigma: 5}}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	size, _ := Size(newImg)
	if size.Height != options.Height || size.Width != options.Width {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	Write("testdata/test_gaussian_out.jpg", newImg)
}

func TestSharpen(t *testing.T) {
	options := Options{Width: 800, Height: 600, Sharpen: Sharpen{Radius: 1, X1: 1.5, Y2: 20, Y3: 50, M1: 1, M2: 2}}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	size, _ := Size(newImg)
	if size.Height != options.Height || size.Width != options.Width {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	Write("testdata/test_sharpen_out.jpg", newImg)
}

func TestExtractWithDefaultAxis(t *testing.T) {
	options := Options{AreaWidth: 200, AreaHeight: 200}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	size, _ := Size(newImg)
	if size.Height != options.AreaHeight || size.Width != options.AreaWidth {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	Write("testdata/test_extract_defaults_out.jpg", newImg)
}

func TestExtractCustomAxis(t *testing.T) {
	options := Options{Top: 100, Left: 100, AreaWidth: 200, AreaHeight: 200}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	size, _ := Size(newImg)
	if size.Height != options.AreaHeight || size.Width != options.AreaWidth {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	Write("testdata/test_extract_custom_axis_out.jpg", newImg)
}

func TestConvert(t *testing.T) {
	width, height := 300, 240
	formats := [3]ImageType{PNG, WEBP, JPEG}

	files := []string{
		"test.jpg",
		"test.png",
		"test.webp",
	}

	for _, file := range files {
		img, err := os.Open("testdata/" + file)
		if err != nil {
			t.Fatal(err)
		}

		buf, err := ioutil.ReadAll(img)
		if err != nil {
			t.Fatal(err)
		}
		img.Close()

		for _, format := range formats {
			options := Options{Width: width, Height: height, Crop: true, Type: format}

			newImg, err := Resize(buf, options)
			if err != nil {
				t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
			}

			if DetermineImageType(newImg) != format {
				t.Fatal("Image is not png")
			}

			size, _ := Size(newImg)
			if size.Height != height || size.Width != width {
				t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
			}
		}
	}
}

func TestResizePngWithTransparency(t *testing.T) {
	width, height := 300, 240

	options := Options{Width: width, Height: height, Crop: true}
	img, err := os.Open("testdata/transparent.png")
	if err != nil {
		t.Fatal(err)
	}
	defer img.Close()

	buf, err := ioutil.ReadAll(img)
	if err != nil {
		t.Fatal(err)
	}

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != PNG {
		t.Fatal("Image is not png")
	}

	size, _ := Size(newImg)
	if size.Height != height || size.Width != width {
		t.Fatal("Invalid image size")
	}

	Write("testdata/transparent_out.png", newImg)
}

func TestRotationAndFlip(t *testing.T) {
	files := []struct {
		Name  string
		Angle Angle
		Flip  bool
	}{
		{"Landscape_1", 0, false},
		{"Landscape_2", 0, true},
		{"Landscape_3", D180, false},
		{"Landscape_4", D180, true},
		{"Landscape_5", D90, true},
		{"Landscape_6", D90, false},
		{"Landscape_7", D270, true},
		{"Landscape_8", D270, false},
		{"Portrait_1", 0, false},
		{"Portrait_2", 0, true},
		{"Portrait_3", D180, false},
		{"Portrait_4", D180, true},
		{"Portrait_5", D90, true},
		{"Portrait_6", D90, false},
		{"Portrait_7", D270, true},
		{"Portrait_8", D270, false},
	}

	for _, file := range files {
		img, err := os.Open(fmt.Sprintf("testdata/exif/%s.jpg", file.Name))
		if err != nil {
			t.Fatal(err)
		}

		buf, err := ioutil.ReadAll(img)
		if err != nil {
			t.Fatal(err)
		}
		img.Close()

		image, _, err := loadImage(buf)
		if err != nil {
			t.Fatal(err)
		}

		angle, flip := calculateRotationAndFlip(image, D0)
		if angle != file.Angle {
			t.Errorf("Rotation for %v expected to be %v. got %v", file.Name, file.Angle, angle)
		}
		if flip != file.Flip {
			t.Errorf("Flip for %v expected to be %v. got %v", file.Name, file.Flip, flip)
		}

		// Visual debugging.
		newImg, err := Resize(buf, Options{})
		if err != nil {
			t.Fatal(err)
		}

		Write(fmt.Sprintf("testdata/exif/%s_out.jpg", file.Name), newImg)
	}
}

func TestIfBothSmartCropOptionsAreIdentical(t *testing.T) {
	if !(VipsMajorVersion >= 8 && VipsMinorVersion > 4) {
		t.Skipf("Skipping this test, libvips doesn't meet version requirement %s > 8.4", VipsVersion)
	}

	benchmarkOptions := Options{Width: 100, Height: 100, Crop: true}
	smartCropOptions := Options{Width: 100, Height: 100, Crop: true, SmartCrop: true}
	gravityOptions := Options{Width: 100, Height: 100, Crop: true, Gravity: GravitySmart}

	testImg, err := os.Open("testdata/northern_cardinal_bird.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer testImg.Close()

	testImgByte, err := ioutil.ReadAll(testImg)
	if err != nil {
		t.Fatal(err)
	}

	scImg, err := Resize(testImgByte, smartCropOptions)
	if err != nil {
		t.Fatal(err)
	}

	gImg, err := Resize(testImgByte, gravityOptions)
	if err != nil {
		t.Fatal(err)
	}

	benchmarkImg, err := Resize(testImgByte, benchmarkOptions)
	if err != nil {
		t.Fatal(err)
	}

	sch, gh, bh := md5.Sum(scImg), md5.Sum(gImg), md5.Sum(benchmarkImg)
	if gh == bh || sch == bh {
		t.Error("Expected both options produce a different result from a standard crop.")
	}

	if sch != gh {
		t.Errorf("Expected both options to result in the same output, %x != %x", sch, gh)
	}
}

func runBenchmarkResize(file string, o Options, b *testing.B) {
	buf, _ := Read(path.Join("testdata", file))

	for n := 0; n < b.N; n++ {
		Resize(buf, o)
	}
}

func BenchmarkRotateJpeg(b *testing.B) {
	options := Options{Rotate: 180}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkResizeLargeJpeg(b *testing.B) {
	options := Options{
		Width:  800,
		Height: 600,
	}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkResizePng(b *testing.B) {
	options := Options{
		Width:  200,
		Height: 200,
	}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkResizeWebp(b *testing.B) {
	options := Options{
		Width:  200,
		Height: 200,
	}
	runBenchmarkResize("test.webp", options, b)
}

func BenchmarkConvertToJpeg(b *testing.B) {
	options := Options{Type: JPEG}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkConvertToPng(b *testing.B) {
	options := Options{Type: PNG}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkConvertToWebp(b *testing.B) {
	options := Options{Type: WEBP}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkCropJpeg(b *testing.B) {
	options := Options{
		Width:  800,
		Height: 600,
	}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkCropPng(b *testing.B) {
	options := Options{
		Width:  800,
		Height: 600,
	}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkCropWebp(b *testing.B) {
	options := Options{
		Width:  800,
		Height: 600,
	}
	runBenchmarkResize("test.webp", options, b)
}

func BenchmarkExtractJpeg(b *testing.B) {
	options := Options{
		Top:        100,
		Left:       50,
		AreaWidth:  600,
		AreaHeight: 480,
	}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkExtractPng(b *testing.B) {
	options := Options{
		Top:        100,
		Left:       50,
		AreaWidth:  600,
		AreaHeight: 480,
	}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkExtractWebp(b *testing.B) {
	options := Options{
		Top:        100,
		Left:       50,
		AreaWidth:  600,
		AreaHeight: 480,
	}
	runBenchmarkResize("test.webp", options, b)
}

func BenchmarkZoomJpeg(b *testing.B) {
	options := Options{Zoom: 1}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkZoomPng(b *testing.B) {
	options := Options{Zoom: 1}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkZoomWebp(b *testing.B) {
	options := Options{Zoom: 1}
	runBenchmarkResize("test.webp", options, b)
}

func BenchmarkWatermarkJpeg(b *testing.B) {
	options := Options{
		Watermark: Watermark{
			Text:       "Chuck Norris (c) 2315",
			Opacity:    0.25,
			Width:      200,
			DPI:        100,
			Margin:     150,
			Font:       "sans bold 12",
			Background: Color{255, 255, 255},
		},
	}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkWatermarkPng(b *testing.B) {
	options := Options{
		Watermark: Watermark{
			Text:       "Chuck Norris (c) 2315",
			Opacity:    0.25,
			Width:      200,
			DPI:        100,
			Margin:     150,
			Font:       "sans bold 12",
			Background: Color{255, 255, 255},
		},
	}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkWatermarkWebp(b *testing.B) {
	options := Options{
		Watermark: Watermark{
			Text:       "Chuck Norris (c) 2315",
			Opacity:    0.25,
			Width:      200,
			DPI:        100,
			Margin:     150,
			Font:       "sans bold 12",
			Background: Color{255, 255, 255},
		},
	}
	runBenchmarkResize("test.webp", options, b)
}

func BenchmarkWatermarkImageJpeg(b *testing.B) {
	watermark := readFile("transparent.png")
	options := Options{
		WatermarkImage: WatermarkImage{
			Buf:     watermark,
			Opacity: 0.25,
			Left:    100,
			Top:     100,
		},
	}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkWatermarkImagePng(b *testing.B) {
	watermark := readFile("transparent.png")
	options := Options{
		WatermarkImage: WatermarkImage{
			Buf:     watermark,
			Opacity: 0.25,
			Left:    100,
			Top:     100,
		},
	}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkWatermarkImageWebp(b *testing.B) {
	watermark := readFile("transparent.png")
	options := Options{
		WatermarkImage: WatermarkImage{
			Buf:     watermark,
			Opacity: 0.25,
			Left:    100,
			Top:     100,
		},
	}
	runBenchmarkResize("test.webp", options, b)
}
