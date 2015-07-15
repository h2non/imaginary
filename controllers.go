package main

import "net/http"

func indexController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorReply(w, "Not found", NOT_FOUND)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("imaginary server " + Version))
}

func formController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlForm()))
}

func imageController(w http.ResponseWriter, r *http.Request, buf []byte, Operation Operation) {
	if len(buf) == 0 {
		ErrorReply(w, "Empty payload", BAD_REQUEST)
		return
	}

	mimeType := http.DetectContentType(buf)
	if IsImageMimeTypeSupported(mimeType) == false {
		ErrorReply(w, "Unsupported media type: "+mimeType, UNSUPPORTED)
		return
	}

	opts := readParams(r)
	if opts.Type != "" && ImageType(opts.Type) == 0 {
		ErrorReply(w, "Unsupported output image format: "+opts.Type, BAD_REQUEST)
		return
	}

	image, err := Operation.Run(buf, opts)
	if err != nil {
		ErrorReply(w, "Error while processing the image: "+err.Error(), BAD_REQUEST)
		return
	}

	w.Header().Set("Content-Type", image.Mime)
	w.Write(image.Body)
}
