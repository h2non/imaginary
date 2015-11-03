package main

import (
	"gopkg.in/h2non/bimg.v0"
	"net/http"
)

func indexController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorReply(w, ErrNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"imaginary\": \"" + Version + "\", \"bimg\": \"" + bimg.Version + "\", \"libvips\": \"" + bimg.VipsVersion + "\" }"))
}

func formController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlForm()))
}

func healthController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func imageControllerDispatcher(o ServerOptions, operation Operation) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var imageSource = MatchSource(req)
		if imageSource == nil {
			ErrorReply(w, ErrMissingImageSource)
			return
		}

		buf, err := imageSource.GetImage(req)
		if err != nil {
			ErrorReply(w, err.(Error))
			return
		}

		if len(buf) == 0 {
			ErrorReply(w, ErrEmptyPayload)
			return
		}

		imageController(w, req, buf, operation)
	}
}

func imageController(w http.ResponseWriter, r *http.Request, buf []byte, Operation Operation) {
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
		ErrorReply(w, NewError("Error while processing the image: "+err.Error(), BAD_REQUEST))
		return
	}

	w.Header().Set("Content-Type", image.Mime)
	w.Write(image.Body)
}
