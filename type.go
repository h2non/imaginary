package main

import (
	"gopkg.in/h2non/bimg.v0"
	"strings"
)

func ExtractImageTypeFromMime(mime string) string {
	mime = strings.Split(mime, " ")[0]
	part := strings.Split(mime, "/")
	if len(part) < 2 {
		return ""
	}
	return strings.ToLower(part[1])
}

func IsImageTypeSupported(mime string) bool {
	format := ExtractImageTypeFromMime(mime)
	return bimg.IsTypeNameSupported(format)
}

func ImageType(mime string) bimg.ImageType {
	format := bimg.UNKNOWN
	switch strings.ToLower(mime) {
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
	}
	return mime
}
