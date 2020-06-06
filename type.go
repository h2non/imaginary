package main

import (
	"strings"

	"github.com/h2non/bimg"
)

// ExtractImageTypeFromMime returns the MIME image type.
func ExtractImageTypeFromMime(mime string) string {
	mime = strings.Split(mime, ";")[0]
	parts := strings.Split(mime, "/")
	if len(parts) < 2 {
		return ""
	}
	name := strings.Split(parts[1], "+")[0]
	return strings.ToLower(name)
}

// IsImageMimeTypeSupported returns true if the image MIME
// type is supported by bimg.
func IsImageMimeTypeSupported(mime string) bool {
	format := ExtractImageTypeFromMime(mime)

	// Some payloads may expose the MIME type for SVG as text/xml
	if format == "xml" {
		format = "svg"
	}

	return bimg.IsTypeNameSupported(format)
}

// ImageType returns the image type based on the given image type alias.
func ImageType(name string) bimg.ImageType {
	switch strings.ToLower(name) {
	case "jpeg":
		return bimg.JPEG
	case "png":
		return bimg.PNG
	case "webp":
		return bimg.WEBP
	case "tiff":
		return bimg.TIFF
	case "gif":
		return bimg.GIF
	case "svg":
		return bimg.SVG
	case "pdf":
		return bimg.PDF
	default:
		return bimg.UNKNOWN
	}
}

// GetImageMimeType returns the MIME type based on the given image type code.
func GetImageMimeType(code bimg.ImageType) string {
	switch code {
	case bimg.PNG:
		return "image/png"
	case bimg.WEBP:
		return "image/webp"
	case bimg.TIFF:
		return "image/tiff"
	case bimg.GIF:
		return "image/gif"
	case bimg.SVG:
		return "image/svg+xml"
	case bimg.PDF:
		return "application/pdf"
	default:
		return "image/jpeg"
	}
}
