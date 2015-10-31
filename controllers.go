package main

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/h2non/bimg.v0"
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

func imageControllerDispatcher(o ServerOptions, operation Operation) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf []byte
		var err error

		if r.Method == "GET" && o.Mount != "" {
			buf, err = readLocalImage(w, r, o.Mount)
		} else {
			buf, err = readPayload(w, r)
		}

		if err != nil {
			return
		}

		if r.Method == "GET" && o.HttpCacheTtl > -1 {
			addCacheHeaders(w, o.HttpCacheTtl)
		}

		imageController(w, r, buf, operation)
	}
}

func addCacheHeaders(w http.ResponseWriter, ttl int) {
	var headerVal string

	ttlDifference := time.Duration(ttl) * time.Second
	expires := time.Now().Add(ttlDifference)

	if ttl == 0 {
		headerVal = "private, no-cache, no-store, must-revalidate"
	} else {
		headerVal = fmt.Sprintf("public, s-maxage: %d, max-age: %d, no-transform", ttl, ttl)
	}

	w.Header().Add("Expires", expires.Format(time.RFC1123))
	w.Header().Add("Cache-Control", headerVal)
}

func imageController(w http.ResponseWriter, r *http.Request, buf []byte, Operation Operation) {
	if len(buf) == 0 {
		ErrorReply(w, ErrEmptyPayload)
		return
	}

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
