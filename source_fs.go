package main

import (
	"io/ioutil"
	"net/http"
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
	return r.Method == http.MethodGet && s.getFileParam(r) != ""
}

func (s *FileSystemImageSource) GetImage(r *http.Request) ([]byte, error) {
	file := s.getFileParam(r)
	if file == "" {
		return nil, ErrMissingParamFile
	}

	file, err := s.buildPath(file)
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

func (s *FileSystemImageSource) getFileParam(r *http.Request) string {
	return r.URL.Query().Get("file")
}

func init() {
	RegisterSource(ImageSourceTypeFileSystem, NewFileSystemImageSource)
}
