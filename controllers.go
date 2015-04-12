package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func indexController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("imaginary server " + Version))
}

const formText = `
<html>
<body>
<h1>Resize</h1>
<form method="POST" action="/resize?width=300&height=200&type=png" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Crop</h1>
<form method="POST" action="/crop" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Flip</h1>
<form method="POST" action="/flip" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Flop</h1>
<form method="POST" action="/flop" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Rotate (180)</h1>
<form method="POST" action="/rotate?rotate=180" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Thumbnail</h1>
<form method="POST" action="/thumbnail?width=100" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Zoom</h1>
<form method="POST" action="/zoom?factor=2&width=300&height=300" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Watermark</h1>
<form method="POST" action="/watermark?text=Hello&font=sans%2014&opacity=0.5" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Convert</h1>
<form method="POST" action="/convert?type=png" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
</body>
</html>
`

func formController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(formText))
}

func infoController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(formText))
}

func mainController(fn Operation) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		imageController(w, r, Operation(fn))
	}
}

func imageController(w http.ResponseWriter, r *http.Request, Operation Operation) {
	if r.Method != "POST" {
		errorResponse(w, "Method not allowed for this endpoint", NOT_ALLOWED)
		return
	}

	var buf []byte
	var err error

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/") {
		err = r.ParseMultipartForm(maxMemory)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		buf, err = getFormPayload(r)
		if err != nil {
			errorResponse(w, "Error while reading the body: "+err.Error(), BAD_REQUEST)
			return
		}
	} else {
		buf, _ = ioutil.ReadAll(r.Body)
	}

	if len(buf) == 0 {
		errorResponse(w, "Empty or invalid body", BAD_REQUEST)
		return
	}

	mimeType := http.DetectContentType(buf)
	if IsImageTypeSupported(mimeType) == false {
		errorResponse(w, "Unsupported media type: "+mimeType, UNSUPPORTED)
		return
	}

	opts, err := parseQueryParams(r)
	if err != nil {
		errorResponse(w, err.Error(), BAD_REQUEST)
		return
	}

	if opts.Type != "" {
		format := ImageType(opts.Type)
		if format == 0 {
			errorResponse(w, "Unsupported image format: "+opts.Type, BAD_REQUEST)
			return
		}
		mimeType = GetImageMimeType(format)
	}

	debug("Options: %#v", opts)
	body, err := Operation.Run(buf, opts)
	if err != nil {
		errorResponse(w, "Error while processing the image: "+err.Error(), BAD_REQUEST)
		return
	}

	w.Header().Set("Content-Type", mimeType)
	w.Write(body)
}

var allowedParams = map[string]string{
	"width":       "int",
	"height":      "int",
	"quality":     "int",
	"top":         "int",
	"left":        "int",
	"compression": "int",
	"rotate":      "int",
	"margin":      "int",
	"factor":      "int",
	"opacity":     "float",
	"noreplicate": "bool",
	"text":        "string",
	"font":        "string",
	"format":      "string",
	"type":        "string",
}

func parseQueryParams(r *http.Request) (ImageOptions, error) {
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
			if param != "" {
				val = true
			}
			break
		}

		params[key] = val
	}

	opts := ImageOptions{
		Width:       params["width"].(int),
		Height:      params["height"].(int),
		Quality:     params["quality"].(int),
		Compression: params["compression"].(int),
		Rotate:      params["rotate"].(int),
		Factor:      params["factor"].(int),
		Text:        params["text"].(string),
		Font:        params["font"].(string),
		Type:        params["type"].(string),
		Opacity:     params["opacity"].(float64),
	}

	return opts, nil
}

func getFormPayload(r *http.Request) ([]byte, error) {
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
