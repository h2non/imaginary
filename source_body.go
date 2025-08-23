package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const formFieldName = "file"
const maxMemory int64 = 1024 * 1024 * 64

const ImageSourceTypeBody ImageSourceType = "payload"

type BodyImageSource struct {
	Config *SourceConfig
}

func NewBodyImageSource(config *SourceConfig) ImageSource {
	return &BodyImageSource{config}
}

func (s *BodyImageSource) Matches(r *http.Request) bool {
	return r.Method == http.MethodPost || r.Method == http.MethodPut
}

func (s *BodyImageSource) GetImage(r *http.Request) ([]byte, error) {
	if isFormBody(r) {
		return readFormBody(r)
	}
	return readRawBody(r)
}

func isFormBody(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/")
}

func readFormBody(r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile(formFieldName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	defaultMaxSize := 10 * 1024 * 1024
	limitedReader := io.LimitReader(file, int64(defaultMaxSize))

	buf, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		err = ErrEmptyBody
	}

	// Check if we actually read the full content (indicating the limit was reached)
	if len(buf) >= defaultMaxSize {
		return nil, fmt.Errorf("form file size exceeds maximum allowed %d bytes", defaultMaxSize)
	}

	return buf, err
}

func readRawBody(r *http.Request) ([]byte, error) {
	defaultMaxSize := 10 * 1024 * 1024

	limitedReader := io.LimitReader(r.Body, int64(defaultMaxSize))
	buf, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	if len(buf) >= defaultMaxSize {
		return nil, fmt.Errorf("request body size exceeds maximum allowed %d bytes", defaultMaxSize)
	}

	return buf, nil
}

func init() {
	RegisterSource(ImageSourceTypeBody, NewBodyImageSource)
}
