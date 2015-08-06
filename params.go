package main

import (
	"gopkg.in/h2non/bimg.v0"
	"math"
	"net/url"
	"strconv"
	"strings"
)

var allowedParams = map[string]string{
	"width":       "int",
	"height":      "int",
	"quality":     "int",
	"top":         "int",
	"left":        "int",
	"areawidth":   "int",
	"areaheight":  "int",
	"compression": "int",
	"rotate":      "int",
	"margin":      "int",
	"factor":      "int",
	"dpi":         "int",
	"textwidth":   "int",
	"opacity":     "float",
	"nocrop":      "bool",
	"noprofile":   "bool",
	"norotation":  "bool",
	"noreplicate": "bool",
	"force":       "bool",
	"text":        "string",
	"font":        "string",
	"type":        "string",
	"color":       "color",
	"colorspace":  "colorspace",
}

func readParams(query url.Values) ImageOptions {
	params := make(map[string]interface{})

	for key, kind := range allowedParams {
		param := query.Get(key)
		params[key] = parseParam(param, kind)
	}

	return mapImageParams(params)
}

func parseParam(param, kind string) interface{} {
	if kind == "int" {
		return parseInt(param)
	}
	if kind == "float" {
		return parseFloat(param)
	}
	if kind == "color" {
		return parseColor(param)
	}
	if kind == "colorspace" {
		return parseColorspace(param)
	}
	if kind == "bool" {
		val, _ := strconv.ParseBool(param)
		return val
	}
	return param
}

func mapImageParams(params map[string]interface{}) ImageOptions {
	return ImageOptions{
		Width:       params["width"].(int),
		Height:      params["height"].(int),
		Top:         params["top"].(int),
		Left:        params["left"].(int),
		AreaWidth:   params["areawidth"].(int),
		AreaHeight:  params["areaheight"].(int),
		DPI:         params["dpi"].(int),
		Quality:     params["quality"].(int),
		TextWidth:   params["textwidth"].(int),
		Compression: params["compression"].(int),
		Rotate:      params["rotate"].(int),
		Factor:      params["factor"].(int),
		Color:       params["color"].([]uint8),
		Text:        params["text"].(string),
		Font:        params["font"].(string),
		Type:        params["type"].(string),
		NoCrop:      params["nocrop"].(bool),
		Force:       params["force"].(bool),
		NoReplicate: params["noreplicate"].(bool),
		NoRotation:  params["norotation"].(bool),
		NoProfile:   params["noprofile"].(bool),
		Opacity:     float32(params["opacity"].(float64)),
		Colorspace:  params["colorspace"].(bimg.Interpretation),
	}
}

func parseColor(val string) []uint8 {
	const max float64 = 255
	buf := []uint8{}
	if val != "" {
		for _, num := range strings.Split(val, ",") {
			n, _ := strconv.ParseUint(strings.Trim(num, " "), 10, 8)
			buf = append(buf, uint8(math.Min(float64(n), max)))
		}
	}
	return buf
}

func parseInt(param string) int {
	return int(math.Floor(parseFloat(param) + 0.5))
}

func parseFloat(param string) float64 {
	val, _ := strconv.ParseFloat(param, 64)
	return math.Abs(val)
}

func parseColorspace(val string) bimg.Interpretation {
	if val == "bw" {
		return bimg.INTERPRETATION_B_W
	}
	return bimg.INTERPRETATION_sRGB
}
