package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	UNAVAILABLE uint8 = iota
	BAD_REQUEST
	NOT_ALLOWED
	UNSUPPORTED
	UNAUTHORIZED
	INTERNAL
	NOT_FOUND
)

var (
	ErrNotFound           = NewError("Not found", NOT_FOUND)
	ErrInvalidApiKey      = NewError("Invalid or missing API key", UNAUTHORIZED)
	ErrMethodNotAllowed   = NewError("Method not allowed", NOT_ALLOWED)
	ErrUnsupportedMedia   = NewError("Unsupported media type", UNSUPPORTED)
	ErrOutputFormat       = NewError("Unsupported output image format", BAD_REQUEST)
	ErrEmptyBody          = NewError("Empty image", BAD_REQUEST)
	ErrMissingParamFile   = NewError("Missing required param: file", BAD_REQUEST)
	ErrInvalidFilePath    = NewError("Invalid file path", BAD_REQUEST)
	ErrInvalidImageURL    = NewError("Invalid image URL", BAD_REQUEST)
	ErrMissingImageSource = NewError("Cannot process the image source", BAD_REQUEST)
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

func (e Error) HttpCode() int {
	if e.Code == BAD_REQUEST {
		return http.StatusBadRequest
	}
	if e.Code == NOT_ALLOWED {
		return http.StatusMethodNotAllowed
	}
	if e.Code == UNSUPPORTED {
		return http.StatusUnsupportedMediaType
	}
	if e.Code == INTERNAL {
		return http.StatusInternalServerError
	}
	if e.Code == UNAUTHORIZED {
		return http.StatusUnauthorized
	}
	if e.Code == NOT_FOUND {
		return http.StatusNotFound
	}
	return http.StatusServiceUnavailable
}

func NewError(err string, code uint8) Error {
	err = strings.Replace(err, "\n", "", -1)
	return Error{err, code}
}

func ErrorReply(w http.ResponseWriter, err Error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HttpCode())
	w.Write(err.JSON())
	return err
}
