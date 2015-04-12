package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const maxMemory int64 = 1024 * 1024 * 64

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
	"noreplicate": "bool",
	"text":        "string",
	"font":        "string",
	"format":      "string",
	"type":        "string",
	"color":       "color",
}

func readBody(r *http.Request) ([]byte, error) {
	var err error
	var buf []byte

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/") {
		err = r.ParseMultipartForm(maxMemory)
		if err != nil {
			return nil, err
		}

		buf, err = readFormPayload(r)
	} else {
		buf, err = ioutil.ReadAll(r.Body)
	}

	return buf, err
}

func readParams(r *http.Request) ImageOptions {
	var val interface{}
	query := r.URL.Query()
	params := make(map[string]interface{})

	for key, kind := range allowedParams {
		key = strings.ToLower(key)
		param := query.Get(key)

		switch kind {
		case "int":
			val, _ = strconv.Atoi(param)
			break
		case "float":
			val, _ = strconv.ParseFloat(param, 64)
			break
		case "string":
			val = param
			break
		case "bool":
			val, _ = strconv.ParseBool(param)
			break
		case "color":
			val = parseColor(param)
			break
		}

		params[key] = val
	}

	return mapImageParams(params)
}

func mapImageParams(params map[string]interface{}) ImageOptions {
	return ImageOptions{
		Width:       params["width"].(int),
		Height:      params["height"].(int),
		Top:         params["top"].(int),
		Left:        params["left"].(int),
		AreaWidth:   params["areawidth"].(int),
		AreaHeight:  params["areawidth"].(int),
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
		Opacity:     params["opacity"].(float64),
		NoReplicate: params["noreplicate"].(bool),
	}
}

func readFormPayload(r *http.Request) ([]byte, error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		return nil, NewError("Empty payload", BAD_REQUEST)
	}

	return buf, err
}

func parseColor(val string) []uint8 {
	buf := []uint8{}
	if val != "" {
		for _, num := range strings.Split(val, ",") {
			n, _ := strconv.ParseUint(num, 10, 8)
			if n < 256 {
				buf = append(buf, uint8(n))
			}
		}
	}
	return buf
}
