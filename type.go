package main

import (
	"gopkg.in/h2non/bimg.v0"
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
	format := bimg.UNKNOWN
	switch strings.ToLower(name) {
	case "jpeg":
		format = bimg.JPEG
		break
	case "png":
		format = bimg.PNG
		break
	case "webp":
		format = bimg.WEBP
		break
	case "tiff":
		format = bimg.TIFF
		break
	}
	return format
}

func GetImageMimeType(code bimg.ImageType) string {
	mime := "image/jpeg"
	switch code {
	case bimg.PNG:
		mime = "image/png"
		break
	case bimg.WEBP:
		mime = "image/webp"
		break
	case bimg.TIFF:
		mime = "image/tiff"
		break
	}
	return mime
}
