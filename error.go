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
	ErrNotFound         = NewError("Not found", NOT_FOUND)
	ErrInvalidApiKey    = NewError("Invalid or missing API key", UNAUTHORIZED)
	ErrMethodNotAllowed = NewError("Method not allowed", NOT_ALLOWED)
	ErrUnsupportedMedia = NewError("Unsupported media type", UNSUPPORTED)
	ErrOutputFormat     = NewError("Unsupported output image format", BAD_REQUEST)
	ErrEmptyPayload     = NewError("Empty payload", BAD_REQUEST)
	ErrMissingParamFile = NewError("Missing required param: file", BAD_REQUEST)
	ErrInvalidFilePath  = NewError("Invalid file path", BAD_REQUEST)
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
	code := http.StatusServiceUnavailable
	switch e.Code {
	case BAD_REQUEST:
		code = http.StatusBadRequest
		break
	case NOT_ALLOWED:
		code = http.StatusMethodNotAllowed
		break
	case UNSUPPORTED:
		code = http.StatusUnsupportedMediaType
		break
	case INTERNAL:
		code = http.StatusInternalServerError
		break
	case UNAUTHORIZED:
		code = http.StatusUnauthorized
		break
	case NOT_FOUND:
		code = http.StatusNotFound
		break
	}
	return code
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
