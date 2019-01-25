package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/h2non/bimg.v1"
)

var ErrUnsupportedValue = errors.New("unsupported value")

// Coercion is the type that type coerces a parameter and defines the appropriate field on ImageOptions
type Coercion func(*ImageOptions, interface{}) error

var paramTypeCoercions = map[string]Coercion{
	"width":       coerceWidth,
	"height":      coerceHeight,
	"quality":     coerceQuality,
	"top":         coerceTop,
	"left":        coerceLeft,
	"areawidth":   coerceAreaWidth,
	"areaheight":  coerceAreaHeight,
	"compression": coerceCompression,
	"rotate":      coerceRotate,
	"margin":      coerceMargin,
	"factor":      coerceFactor,
	"dpi":         coerceDPI,
	"textwidth":   coerceTextWidth,
	"opacity":     coerceOpacity,
	"flip":        coerceFlip,
	"flop":        coerceFlop,
	"nocrop":      coerceNoCrop,
	"noprofile":   coerceNoProfile,
	"norotation":  coerceNoRotation,
	"noreplicate": coerceNoReplicate,
	"force":       coerceForce,
	"embed":       coerceEmbed,
	"stripmeta":   coerceStripMeta,
	"text":        coerceText,
	"image":       coerceImage,
	"font":        coerceFont,
	"type":        coerceImageType,
	"color":       coerceColor,
	"colorspace":  coerceColorSpace,
	"gravity":     coerceGravity,
	"background":  coerceBackground,
	"extend":      coerceExtend,
	"sigma":       coerceSigma,
	"minampl":     coerceMinAmpl,
	"operations":  coerceOperations,
}

func coerceTypeInt(param interface{}) (int, error) {
	if v, ok := param.(int); ok {
		return v, nil
	}

	if v, ok := param.(float64); ok {
		return int(v), nil
	}

	if v, ok := param.(string); ok {
		return parseInt(v)
	}

	return 0, ErrUnsupportedValue
}

func coerceTypeFloat(param interface{}) (float64, error) {
	if v, ok := param.(float64); ok {
		return v, nil
	}

	if v, ok := param.(int); ok {
		return float64(v), nil
	}

	if v, ok := param.(string); ok {
		result, err := parseFloat(v)
		if err != nil {
			return 0, ErrUnsupportedValue
		}

		return result, nil
	}

	return 0, ErrUnsupportedValue
}

func coerceTypeBool(param interface{}) (bool, error) {
	if v, ok := param.(bool); ok {
		return v, nil
	}

	if v, ok := param.(string); ok {
		result, err := parseBool(v)
		if err != nil {
			return false, ErrUnsupportedValue
		}

		return result, nil
	}

	return false, ErrUnsupportedValue
}

func coerceTypeString(param interface{}) (string, error) {
	if v, ok := param.(string); ok {
		return v, nil
	}

	return "", ErrUnsupportedValue
}

func coerceHeight(io *ImageOptions, param interface{}) (err error) {
	io.Height, err = coerceTypeInt(param)
	return err
}

func coerceWidth(io *ImageOptions, param interface{}) (err error) {
	io.Width, err = coerceTypeInt(param)
	return err
}

func coerceQuality(io *ImageOptions, param interface{}) (err error) {
	io.Quality, err = coerceTypeInt(param)
	return err
}

func coerceTop(io *ImageOptions, param interface{}) (err error) {
	io.Top, err = coerceTypeInt(param)
	return err
}

func coerceLeft(io *ImageOptions, param interface{}) (err error) {
	io.Left, err = coerceTypeInt(param)
	return err
}

func coerceAreaWidth(io *ImageOptions, param interface{}) (err error) {
	io.AreaWidth, err = coerceTypeInt(param)
	return err
}

func coerceAreaHeight(io *ImageOptions, param interface{}) (err error) {
	io.AreaHeight, err = coerceTypeInt(param)
	return err
}

func coerceCompression(io *ImageOptions, param interface{}) (err error) {
	io.Compression, err = coerceTypeInt(param)
	return err
}

func coerceRotate(io *ImageOptions, param interface{}) (err error) {
	io.Rotate, err = coerceTypeInt(param)
	return err
}

func coerceMargin(io *ImageOptions, param interface{}) (err error) {
	io.Margin, err = coerceTypeInt(param)
	return err
}

func coerceFactor(io *ImageOptions, param interface{}) (err error) {
	io.Factor, err = coerceTypeInt(param)
	return err
}

func coerceDPI(io *ImageOptions, param interface{}) (err error) {
	io.DPI, err = coerceTypeInt(param)
	return err
}

func coerceTextWidth(io *ImageOptions, param interface{}) (err error) {
	io.TextWidth, err = coerceTypeInt(param)
	return err
}

func coerceOpacity(io *ImageOptions, param interface{}) (err error) {
	v, err := coerceTypeFloat(param)
	io.Opacity = float32(v)
	return err
}

func coerceFlip(io *ImageOptions, param interface{}) (err error) {
	io.Flip, err = coerceTypeBool(param)
	io.IsDefinedField.Flip = true
	return err
}

func coerceFlop(io *ImageOptions, param interface{}) (err error) {
	io.Flop, err = coerceTypeBool(param)
	io.IsDefinedField.Flop = true
	return err
}

func coerceNoCrop(io *ImageOptions, param interface{}) (err error) {
	io.NoCrop, err = coerceTypeBool(param)
	io.IsDefinedField.NoCrop = true
	return err
}

func coerceNoProfile(io *ImageOptions, param interface{}) (err error) {
	io.NoProfile, err = coerceTypeBool(param)
	io.IsDefinedField.NoProfile = true
	return err
}

func coerceNoRotation(io *ImageOptions, param interface{}) (err error) {
	io.NoRotation, err = coerceTypeBool(param)
	io.IsDefinedField.NoRotation = true
	return err
}

func coerceNoReplicate(io *ImageOptions, param interface{}) (err error) {
	io.NoReplicate, err = coerceTypeBool(param)
	io.IsDefinedField.NoReplicate = true
	return err
}

func coerceForce(io *ImageOptions, param interface{}) (err error) {
	io.Force, err = coerceTypeBool(param)
	io.IsDefinedField.Force = true
	return err
}

func coerceEmbed(io *ImageOptions, param interface{}) (err error) {
	io.Embed, err = coerceTypeBool(param)
	io.IsDefinedField.Embed = true
	return err
}

func coerceStripMeta(io *ImageOptions, param interface{}) (err error) {
	io.StripMetadata, err = coerceTypeBool(param)
	io.IsDefinedField.StripMetadata = true
	return err
}

func coerceText(io *ImageOptions, param interface{}) (err error) {
	io.Text, err = coerceTypeString(param)
	return err
}

func coerceImage(io *ImageOptions, param interface{}) (err error) {
	io.Image, err = coerceTypeString(param)
	return err
}

func coerceFont(io *ImageOptions, param interface{}) (err error) {
	io.Font, err = coerceTypeString(param)
	return err
}

func coerceImageType(io *ImageOptions, param interface{}) (err error) {
	io.Type, err = coerceTypeString(param)
	return err
}

func coerceColor(io *ImageOptions, param interface{}) error {
	if v, ok := param.(string); ok {
		io.Color = parseColor(v)
		return nil
	}

	return ErrUnsupportedValue
}

func coerceColorSpace(io *ImageOptions, param interface{}) error {
	if v, ok := param.(string); ok {
		io.Colorspace = parseColorspace(v)
		return nil
	}

	return ErrUnsupportedValue
}

func coerceGravity(io *ImageOptions, param interface{}) error {
	if v, ok := param.(string); ok {
		io.Gravity = parseGravity(v)
		return nil
	}

	return ErrUnsupportedValue
}

func coerceBackground(io *ImageOptions, param interface{}) error {
	if v, ok := param.(string); ok {
		io.Background = parseColor(v)
		return nil
	}

	return ErrUnsupportedValue
}

func coerceExtend(io *ImageOptions, param interface{}) error {
	if v, ok := param.(string); ok {
		io.Extend = parseExtendMode(v)
		return nil
	}

	return ErrUnsupportedValue
}

func coerceSigma(io *ImageOptions, param interface{}) (err error) {
	io.Sigma, err = coerceTypeFloat(param)
	return err
}

func coerceMinAmpl(io *ImageOptions, param interface{}) (err error) {
	io.MinAmpl, err = coerceTypeFloat(param)
	return err
}

func coerceOperations(io *ImageOptions, param interface{}) (err error) {
	if v, ok := param.(string); ok {
		ops, err := parseJSONOperations(v)
		if err == nil {
			io.Operations = ops
		}

		return err
	}

	return ErrUnsupportedValue
}

func buildParamsFromOperation(op PipelineOperation) (ImageOptions, error) {

	var options ImageOptions

	for key, value := range op.Params {
		fn, ok := paramTypeCoercions[key]
		if !ok {
			continue
		}

		err := fn(&options, value)
		if err != nil {
			return ImageOptions{}, fmt.Errorf(`error while processing parameter "%s" with value %q, error: %s`, key, value, err)
		}
	}

	return options, nil
}

// buildParamsFromQuery builds the ImageOptions type from untyped parameters
func buildParamsFromQuery(query url.Values) (ImageOptions, error) {
	var options ImageOptions

	// Extract only known parameters
	for key := range query {
		fn, ok := paramTypeCoercions[key]
		if !ok {
			continue
		}

		value := query.Get(key)
		err := fn(&options, value)
		if err != nil {
			return ImageOptions{}, fmt.Errorf(`error while processing parameter "%s" with value %q, error: %s`, key, value, err)
		}
	}

	return options, nil
}

func parseBool(val string) (bool, error) {
	if val == "" {
		return false, nil
	}

	return strconv.ParseBool(val)
}

func parseInt(param string) (int, error) {
	if param == "" {
		return 0, nil
	}

	f, err := parseFloat(param)
	return int(math.Floor(f + 0.5)), err
}

func parseFloat(param string) (float64, error) {
	if param == "" {
		return 0.0, nil
	}

	val, err := strconv.ParseFloat(param, 64)
	return math.Abs(val), err
}

func parseColorspace(val string) bimg.Interpretation {
	if val == "bw" {
		return bimg.InterpretationBW
	}
	return bimg.InterpretationSRGB
}

func parseColor(val string) []uint8 {
	const max float64 = 255
	var buf []uint8
	if val != "" {
		for _, num := range strings.Split(val, ",") {
			n, _ := strconv.ParseUint(strings.Trim(num, " "), 10, 8)
			buf = append(buf, uint8(math.Min(float64(n), max)))
		}
	}
	return buf
}

func parseJSONOperations(data string) (PipelineOperations, error) {
	var operations PipelineOperations

	// Fewer than 2 characters cannot be valid JSON. We assume empty operation.
	if len(data) < 2 {
		return operations, nil
	}

	d := json.NewDecoder(strings.NewReader(data))
	d.DisallowUnknownFields()

	err := d.Decode(&operations)
	return operations, err
}

func parseExtendMode(val string) bimg.Extend {
	val = strings.TrimSpace(strings.ToLower(val))
	if val == "white" {
		return bimg.ExtendWhite
	}
	if val == "copy" {
		return bimg.ExtendCopy
	}
	if val == "mirror" {
		return bimg.ExtendMirror
	}
	if val == "background" {
		return bimg.ExtendBackground
	}
	return bimg.ExtendBlack
}

func parseGravity(val string) bimg.Gravity {
	var m = map[string]bimg.Gravity{
		"south": bimg.GravitySouth,
		"north": bimg.GravityNorth,
		"east":  bimg.GravityEast,
		"west":  bimg.GravityWest,
		"smart": bimg.GravitySmart,
	}

	val = strings.TrimSpace(strings.ToLower(val))
	if g, ok := m[val]; ok {
		return g
	}

	return bimg.GravityCentre
}
