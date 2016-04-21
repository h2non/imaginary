package main

import (
	"gopkg.in/h2non/bimg.v1"
	"strings"
)

func ExtractImageTypeFromMime(mime string) string {
	mime = strings.Split(mime, ";")[0]
	part := strings.Split(mime, "/")
	if len(part) < 2 {
		return ""
	}
	return strings.ToLower(part[1])
}

func IsImageMimeTypeSupported(mime string) bool {
	format := ExtractImageTypeFromMime(mime)
	return bimg.IsTypeNameSupported(format)
}

func ImageType(name string) bimg.ImageType {
	ext := strings.ToLower(name)
	if ext == "jpeg" {
		return bimg.JPEG
	}
	if ext == "png" {
		return bimg.PNG
	}
	if ext == "webp" {
		return bimg.WEBP
	}
	if ext == "tiff" {
		return bimg.TIFF
	}
	return bimg.UNKNOWN
}

func GetImageMimeType(code bimg.ImageType) string {
	if code == bimg.PNG {
		return "image/png"
	}
	if code == bimg.WEBP {
		return "image/webp"
	}
	if code == bimg.TIFF {
		return "image/tiff"
	}
	return "image/jpeg"
}
