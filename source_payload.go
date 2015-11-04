package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const maxMemory int64 = 1024 * 1024 * 32

const ImageSourceTypePayload ImageSourceType = "payload"

type PayloadImageSource struct {
	Config *SourceConfig
}

func NewPayloadImageSource(config *SourceConfig) ImageSource {
	return &PayloadImageSource{config}
}

func (s *PayloadImageSource) Matches(r *http.Request) bool {
	return r.Method == "POST"
}

func (s *PayloadImageSource) GetImage(r *http.Request) ([]byte, error) {
	if isFormPayload(r) {
		return readFormPayload(r)
	}
	return readRawPayload(r)
}

func isFormPayload(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/")
}

func readFormPayload(r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if len(buf) == 0 {
		err = ErrEmptyPayload
	}

	return buf, err
}

func readRawPayload(r *http.Request) ([]byte, error) {
	return ioutil.ReadAll(r.Body)
}

func init() {
	RegisterSource(ImageSourceTypePayload, NewPayloadImageSource)
}
