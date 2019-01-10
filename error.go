package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	bimg "gopkg.in/h2non/bimg.v1"
)

const (
	_ uint8 = iota
	BadRequest
	NotAllowed
	Unsupported
	Unauthorized
	InternalError
	NotFound
	NotImplemented
	Forbidden
	NotAcceptable
)

var (
	ErrNotFound             = NewError("not found", NotFound)
	ErrInvalidAPIKey        = NewError("invalid or missing API key", Unauthorized)
	ErrMethodNotAllowed     = NewError("method not allowed", NotAllowed)
	ErrUnsupportedMedia     = NewError("unsupported media type", Unsupported)
	ErrOutputFormat         = NewError("unsupported output image format", BadRequest)
	ErrEmptyBody            = NewError("empty image", BadRequest)
	ErrMissingParamFile     = NewError("missing required param: file", BadRequest)
	ErrInvalidFilePath      = NewError("invalid file path", BadRequest)
	ErrInvalidImageURL      = NewError("invalid image URL", BadRequest)
	ErrMissingImageSource   = NewError("cannot process the image due to missing or invalid params", BadRequest)
	ErrNotImplemented       = NewError("not implemented endpoint", NotImplemented)
	ErrInvalidURLSignature  = NewError("invalid URL signature", BadRequest)
	ErrURLSignatureMismatch = NewError("URL signature mismatch", Forbidden)
)

type Error struct {
	Message string `json:"message,omitempty"`
	Code    uint8  `json:"code"`
}

func (e Error) JSON() []byte {
	buf, _ := json.Marshal(e)
	return buf
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) HTTPCode() int {
	var codes = map[uint8]int{
		BadRequest:     http.StatusBadRequest,
		NotAllowed:     http.StatusMethodNotAllowed,
		Unsupported:    http.StatusUnsupportedMediaType,
		InternalError:  http.StatusInternalServerError,
		Unauthorized:   http.StatusUnauthorized,
		NotFound:       http.StatusNotFound,
		NotImplemented: http.StatusNotImplemented,
		Forbidden:      http.StatusForbidden,
		NotAcceptable:  http.StatusNotAcceptable,
	}

	if v, ok := codes[e.Code]; ok {
		return v
	}

	return http.StatusServiceUnavailable
}

func NewError(err string, code uint8) Error {
	err = strings.Replace(err, "\n", "", -1)
	return Error{err, code}
}

func replyWithPlaceholder(req *http.Request, w http.ResponseWriter, err Error, o ServerOptions) error {
	// Resize placeholder to expected output
	buf, _err := bimg.Resize(o.PlaceholderImage, bimg.Options{
		Force:   true,
		Crop:    true,
		Enlarge: true,
		Width:   parseInt(req.URL.Query().Get("width")),
		Height:  parseInt(req.URL.Query().Get("height")),
		Type:    ImageType(req.URL.Query().Get("type")),
	})

	if _err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\", \"code\": %d}", _err.Error(), BadRequest)))
		return _err
	}

	// Use final response body image
	image := buf

	// Placeholder image response
	w.Header().Set("Content-Type", GetImageMimeType(bimg.DetermineImageType(image)))
	w.Header().Set("Error", string(err.JSON()))
	w.WriteHeader(err.HTTPCode())
	_, _ = w.Write(image)

	return err
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
