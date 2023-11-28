package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const ImageSourceTypeFileSystem ImageSourceType = "fs"

type FileSystemImageSource struct {
	Config *SourceConfig
}

func NewFileSystemImageSource(config *SourceConfig) ImageSource {
	return &FileSystemImageSource{config}
}

func (s *FileSystemImageSource) Matches(r *http.Request) bool {
	file, err := s.getFileParam(r)
	if err != nil {
		return false
	}
	return r.Method == http.MethodGet && file != ""
}

func (s *FileSystemImageSource) GetImage(r *http.Request) ([]byte, error) {
	file, err := s.getFileParam(r)
	if err != nil {
		return nil, err
	}

	if file == "" {
		return nil, ErrMissingParamFile
	}

	file, err = s.buildPath(file)
	if err != nil {
		return nil, err
	}

	return s.read(file)
}

func (s *FileSystemImageSource) buildPath(file string) (string, error) {
	file = path.Clean(path.Join(s.Config.MountPath, file))
	if !strings.HasPrefix(file, s.Config.MountPath) {
		return "", ErrInvalidFilePath
	}
	return file, nil
}

func (s *FileSystemImageSource) read(file string) ([]byte, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, ErrInvalidFilePath
	}
	return buf, nil
}

func (s *FileSystemImageSource) getFileParam(r *http.Request) (string, error) {
	unescaped, err := url.QueryUnescape(r.URL.Query().Get("file"))
	if err != nil{
		return "", fmt.Errorf("failed to unescape file param: %w", err)
	}

	return unescaped, nil
}

func init() {
	RegisterSource(ImageSourceTypeFileSystem, NewFileSystemImageSource)
}
