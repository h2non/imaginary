package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/h2non/bimg.v1"
	"gopkg.in/h2non/filetype.v0"
)

func indexController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorReply(r, w, ErrNotFound, ServerOptions{})
		return
	}

	body, _ := json.Marshal(Versions{
		Version,
		bimg.Version,
		bimg.VipsVersion,
	})
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

func healthController(w http.ResponseWriter, r *http.Request) {
	health := GetHealthStats()
	body, _ := json.Marshal(health)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

func imageController(o ServerOptions, operation Operation) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var imageSource = MatchSource(req)
		if imageSource == nil {
			ErrorReply(req, w, ErrMissingImageSource, o)
			return
		}

		buf, err := imageSource.GetImage(req)
		if err != nil {
			ErrorReply(req, w, NewError(err.Error(), BadRequest), o)
			return
		}

		if len(buf) == 0 {
			ErrorReply(req, w, ErrEmptyBody, o)
			return
		}

		imageHandler(w, req, buf, operation, o)
	}
}

func determineAcceptMimeType(accept string) string {
	for _, v := range strings.Split(accept, ",") {
		mediaType, _, _ := mime.ParseMediaType(v)
		switch mediaType {
		case "image/webp":
			return "webp"
		case "image/png":
			return "png"
		case "image/jpeg":
			return "jpeg"
		}
	}

	return ""
}

func imageHandler(w http.ResponseWriter, r *http.Request, buf []byte, operation Operation, o ServerOptions) {
	// Infer the body MIME type via mime sniff algorithm
	mimeType := http.DetectContentType(buf)

	// If cannot infer the type, infer it via magic numbers
	if mimeType == "application/octet-stream" {
		kind, err := filetype.Get(buf)
		if err == nil && kind.MIME.Value != "" {
			mimeType = kind.MIME.Value
		}
	}

	// Infer text/plain responses as potential SVG image
	if strings.Contains(mimeType, "text/plain") && len(buf) > 8 {
		if bimg.IsSVGImage(buf) {
			mimeType = "image/svg+xml"
		}
	}

	// Finally check if image MIME type is supported
	if !IsImageMimeTypeSupported(mimeType) {
		ErrorReply(r, w, ErrUnsupportedMedia, o)
		return
	}

	opts, err := buildParamsFromQuery(r.URL.Query())
	if err != nil {
		ErrorReply(r, w, NewError("Error while processing parameters, "+err.Error(), BadRequest), o)
		return
	}

	vary := ""
	if opts.Type == "auto" {
		opts.Type = determineAcceptMimeType(r.Header.Get("Accept"))
		vary = "Accept" // Ensure caches behave correctly for negotiated content
	} else if opts.Type != "" && ImageType(opts.Type) == 0 {
		ErrorReply(r, w, ErrOutputFormat, o)
		return
	}

	image, err := operation.Run(buf, opts)
	if err != nil {
		ErrorReply(r, w, NewError("Error while processing the image: "+err.Error(), BadRequest), o)
		return
	}

	// Expose Content-Length response header
	w.Header().Set("Content-Length", strconv.Itoa(len(image.Body)))
	w.Header().Set("Content-Type", image.Mime)
	if vary != "" {
		w.Header().Set("Vary", vary)
	}
	_, _ = w.Write(image.Body)
}

func formController(w http.ResponseWriter, r *http.Request) {
	operations := []struct {
		name   string
		method string
		args   string
	}{
		{"Resize", "resize", "width=300&height=200&type=jpeg"},
		{"Force resize", "resize", "width=300&height=200&force=true"},
		{"Crop", "crop", "width=300&quality=95"},
		{"SmartCrop", "crop", "width=300&height=260&quality=95&gravity=smart"},
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
		{"Gaussian blur", "blur", "sigma=15.0&minampl=0.2"},
		{"Pipeline (image reduction via multiple transformations)", "pipeline", "operations=%5B%7B%22operation%22:%20%22crop%22,%20%22params%22:%20%7B%22width%22:%20300,%20%22height%22:%20260%7D%7D,%20%7B%22operation%22:%20%22convert%22,%20%22params%22:%20%7B%22type%22:%20%22webp%22%7D%7D%5D"},
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
	_, _ = w.Write([]byte(html))
}
