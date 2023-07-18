package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/h2non/bimg"
)

var (
	ErrNotFound             = NewError("Not found", http.StatusNotFound)
	ErrInvalidAPIKey        = NewError("Invalid or missing API key", http.StatusUnauthorized)
	ErrMethodNotAllowed     = NewError("HTTP method not allowed. Try with a POST or GET method (-enable-url-source flag must be defined)", http.StatusMethodNotAllowed)
	ErrGetMethodNotAllowed  = NewError("GET method not allowed. Make sure remote URL source is enabled by using the flag: -enable-url-source", http.StatusMethodNotAllowed)
	ErrUnsupportedMedia     = NewError("Unsupported media type", http.StatusNotAcceptable)
	ErrOutputFormat         = NewError("Unsupported output image format", http.StatusBadRequest)
	ErrEmptyBody            = NewError("Empty or unreadable image", http.StatusBadRequest)
	ErrMissingParamFile     = NewError("Missing required param: file", http.StatusBadRequest)
	ErrInvalidFilePath      = NewError("Invalid file path", http.StatusBadRequest)
	ErrInvalidImageURL      = NewError("Invalid image URL", http.StatusBadRequest)
	ErrMissingImageSource   = NewError("Cannot process the image due to missing or invalid params", http.StatusBadRequest)
	ErrNotImplemented       = NewError("Not implemented endpoint", http.StatusNotImplemented)
	ErrInvalidURLSignature  = NewError("Invalid URL signature", http.StatusBadRequest)
	ErrURLSignatureMismatch = NewError("URL signature mismatch", http.StatusForbidden)
	ErrResolutionTooBig     = NewError("Image resolution is too big", http.StatusUnprocessableEntity)
)

type Error struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"status"`
}

func (e Error) JSON() []byte {
	buf, _ := json.Marshal(e)
	return buf
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) HTTPCode() int {
	if e.Code >= 400 && e.Code <= 511 {
		return e.Code
	}
	return http.StatusServiceUnavailable
}

func NewError(err string, code int) Error {
	err = strings.ReplaceAll(err, "\n", "")
	return Error{Message: err, Code: code}
}

func sendErrorResponse(w http.ResponseWriter, httpStatusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	_, _ = w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\", \"status\": %d}", err.Error(), httpStatusCode)))
}

func replyWithPlaceholder(req *http.Request, w http.ResponseWriter, errCaller Error, o ServerOptions) error {
	var err error
	bimgOptions := bimg.Options{
		Force:   true,
		Crop:    true,
		Enlarge: true,
		Type:    ImageType(req.URL.Query().Get("type")),
	}

	bimgOptions.Width, err = parseInt(req.URL.Query().Get("width"))
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err)
		return err
	}

	bimgOptions.Height, err = parseInt(req.URL.Query().Get("height"))
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err)
		return err
	}

	// Resize placeholder to expected output
	buf, err := bimg.Resize(o.PlaceholderImage, bimgOptions)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err)
		return err
	}

	// Use final response body image
	image := buf

	// Placeholder image response
	w.Header().Set("Content-Type", GetImageMimeType(bimg.DetermineImageType(image)))
	w.Header().Set("Error", string(errCaller.JSON()))
	if o.PlaceholderStatus != 0 {
		w.WriteHeader(o.PlaceholderStatus)
	} else {
		w.WriteHeader(errCaller.HTTPCode())
	}
	_, _ = w.Write(image)

	return errCaller
}

func ErrorReply(req *http.Request, w http.ResponseWriter, err Error, o ServerOptions) {
	// Reply with placeholder if required
	if o.EnablePlaceholder || o.Placeholder != "" {
		_ = replyWithPlaceholder(req, w, err, o)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPCode())
	_, _ = w.Write(err.JSON())
}
