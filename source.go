package main

import (
	"fmt"
	"net/http"
	"os"
)

type ImageSourceType string
type ImageSourceFactoryFunction func(*SourceConfig) ImageSource

type SourceConfig struct {
	Type      ImageSourceType
	Directory string
}

var (
	imageSourceTypeToFactoryFunctionMap = make(map[ImageSourceType]ImageSourceFactoryFunction)
)

type ImageSource interface {
	GetImage(*http.Request) ([]byte, error)
}

func RegisterSource(sourceType ImageSourceType, factory ImageSourceFactoryFunction) {
	imageSourceTypeToFactoryFunctionMap[sourceType] = factory
}

func NewImageSourceWithConfig(config *SourceConfig) ImageSource {
	factory := imageSourceTypeToFactoryFunctionMap[config.Type]
	if factory == nil {
		fmt.Fprintf(os.Stderr, "Unknown image source type: %s\n", config.Type)
		os.Exit(1)
	}
	return factory(config)
}
