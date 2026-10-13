package main

import (
	"net/http"
	"net/url"
)

type ImageSourceType string
type ImageSourceFactoryFunction func(*SourceConfig) ImageSource

type SourceConfig struct {
	AuthForwarding bool
	Authorization  string
	MountPath      string
	TmpDir         string
	SpoolThreshold int
	Type           ImageSourceType
	ForwardHeaders []string
	AllowedOrigins []*url.URL
	MaxAllowedSize int
}

var imageSourceMap = make(map[ImageSourceType]ImageSource)
var imageSourceFactoryMap = make(map[ImageSourceType]ImageSourceFactoryFunction)

type ImageSource interface {
	Matches(*http.Request) bool
	GetImage(*http.Request) ([]byte, error)
}

// ReleasingImageSource is an ImageSource whose image may be a mapped file. The buffer
// stays valid until release is called, which the caller does once it is done with it.
type ReleasingImageSource interface {
	GetImageReleasing(*http.Request) (buf []byte, release func(), err error)
}

// copyAndRelease turns the buffer of a ReleasingImageSource into one the caller owns.
func copyAndRelease(buf []byte, release func(), err error) ([]byte, error) {
	defer release()
	if err != nil {
		return nil, err
	}
	return append([]byte(nil), buf...), nil
}

func RegisterSource(sourceType ImageSourceType, factory ImageSourceFactoryFunction) {
	imageSourceFactoryMap[sourceType] = factory
}

func LoadSources(o ServerOptions) {
	for name, factory := range imageSourceFactoryMap {
		imageSourceMap[name] = factory(&SourceConfig{
			Type:           name,
			MountPath:      o.Mount,
			TmpDir:         o.TmpDir,
			SpoolThreshold: o.SpoolThreshold,
			AuthForwarding: o.AuthForwarding,
			Authorization:  o.Authorization,
			AllowedOrigins: o.AllowedOrigins,
			MaxAllowedSize: o.MaxAllowedSize,
			ForwardHeaders: o.ForwardHeaders,
		})
	}
}

func MatchSource(req *http.Request) ImageSource {
	for _, source := range imageSourceMap {
		if source.Matches(req) {
			return source
		}
	}
	return nil
}
