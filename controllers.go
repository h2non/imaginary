package main

import (
	"net/http"
)

func indexController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("imaginary server " + Version))
}

func imageController(w http.ResponseWriter, r *http.Request, Operation Operation) {
	if r.Method != "POST" {
		errorResponse(w, "Method not allowed for this endpoint", NOT_ALLOWED)
		return
	}

	buf, err := readBody(r)
	if err != nil {
		errorResponse(w, "Cannot read the body: "+err.Error(), BAD_REQUEST)
		return
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

	opts := readParams(r)

	if opts.Type != "" && ImageType(opts.Type) == 0 {
		errorResponse(w, "Unsupported conversion image format: "+opts.Type, BAD_REQUEST)
		return
	}

	image, err := Operation.Run(buf, opts)
	if err != nil {
		errorResponse(w, "Error while processing the image: "+err.Error(), BAD_REQUEST)
		return
	}

	w.Header().Set("Content-Type", image.Mime)
	w.Write(image.Body)
}
