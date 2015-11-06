package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/h2non/bimg.v0"
	"net/http"
)

func indexController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorReply(w, ErrNotFound)
		return
	}

	versions := struct {
		ImaginaryVersion string `json:"imaginary"`
		BimgVersion      string `json:"bimg"`
		VipsVersion      string `json:"libvips"`
	}{Version, bimg.Version, bimg.VipsVersion}

	body, _ := json.Marshal(versions)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func healthController(w http.ResponseWriter, r *http.Request) {
	health := GetHealthStats()
	body, _ := json.Marshal(health)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func imageController(o ServerOptions, operation Operation) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var imageSource = MatchSource(req)
		if imageSource == nil {
			ErrorReply(w, ErrMissingImageSource)
			return
		}

		buf, err := imageSource.GetImage(req)
		if err != nil {
			ErrorReply(w, NewError(err.Error(), BadRequest))
			return
		}

		if len(buf) == 0 {
			ErrorReply(w, ErrEmptyBody)
			return
		}

		imageHandler(w, req, buf, operation)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request, buf []byte, Operation Operation) {
	mimeType := http.DetectContentType(buf)
	if IsImageMimeTypeSupported(mimeType) == false {
		ErrorReply(w, ErrUnsupportedMedia)
		return
	}

	opts := readParams(r.URL.Query())
	if opts.Type != "" && ImageType(opts.Type) == 0 {
		ErrorReply(w, ErrOutputFormat)
		return
	}

	image, err := Operation.Run(buf, opts)
	if err != nil {
		ErrorReply(w, NewError("Error while processing the image: "+err.Error(), BadRequest))
		return
	}

	w.Header().Set("Content-Type", image.Mime)
	w.Write(image.Body)
}

func formController(w http.ResponseWriter, r *http.Request) {
	operations := []struct {
		name   string
		method string
		args   string
	}{
		{"Resize", "resize", "width=300&height=200&type=png"},
		{"Force resize", "resize", "width=300&height=200&force=true"},
		{"Crop", "crop", "width=562&height=562&quality=95"},
		{"Extract", "extract", "top=100&left=100&areawidth=300&areaheight=150"},
		{"Enlarge", "enlarge", "width=1440&height=900&quality=95"},
		{"Rotate", "rotate", "rotate=180"},
		{"Flip", "flip", ""},
		{"Flop", "flop", ""},
		{"Thumbnail", "thumbnail", "width=100"},
		{"Zoom", "zoom", "factor=2&areawidth=300&top=80&left=80"},
		{"Color space (black&white)", "resize", "width=400&height=300&colorspace=bw"},
		{"Add watermark", "watermark", "textwidth=100&text=Hello&font=sans%2012&opacity=0.5&color=255,200,50"},
		{"Convert format", "convert", "type=png"},
		{"Image metadata", "info", ""},
	}

	html := "<html><body>"

	for _, form := range operations {
		html += fmt.Sprintf(`
    <h1>%s</h1>
    <form method="POST" action="/%s?%s" enctype="multipart/form-data">
      <input type="file" name="file" />
      <input type="submit" value="Upload" />
    </form>`, form.name, form.method, form.args)
	}

	html += "</body></html>"

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
