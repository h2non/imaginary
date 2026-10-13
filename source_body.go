package main

import (
	"io/ioutil"
	"net/http"
	"os"
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
	return copyAndRelease(s.GetImageReleasing(r))
}

func (s *BodyImageSource) GetImageReleasing(r *http.Request) ([]byte, func(), error) {
	if isFormBody(r) {
		return readFormBody(r)
	}
	return readSpooled(r.Body, s.Config.TmpDir, int64(s.Config.SpoolThreshold))
}

func isFormBody(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/")
}

func readFormBody(r *http.Request) ([]byte, func(), error) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, noRelease, err
	}

	file, _, err := r.FormFile(formFieldName)
	if err != nil {
		return nil, noRelease, err
	}
	defer file.Close()

	// A part larger than maxMemory is already a temporary file
	if f, ok := file.(*os.File); ok {
		buf, release, err := mapFile(f)
		if err == nil && len(buf) == 0 {
			release()
			return nil, noRelease, ErrEmptyBody
		}
		return buf, release, err
	}

	buf, err := ioutil.ReadAll(file)
	if len(buf) == 0 {
		err = ErrEmptyBody
	}

	return buf, noRelease, err
}

func init() {
	RegisterSource(ImageSourceTypeBody, NewBodyImageSource)
}
