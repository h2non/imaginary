package main

import (
	"strings"

	"gopkg.in/h2non/bimg.v1"
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
	if ext == "gif" {
		return bimg.GIF
	}
	if ext == "svg" {
		return bimg.SVG
	}
	if ext == "pdf" {
		return bimg.PDF
	}
	return bimg.UNKNOWN
}

// GetImageMimeType returns the MIME type based on the given image type code.
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
	if code == bimg.GIF {
		return "image/gif"
	}
	if code == bimg.SVG {
		return "image/svg+xml"
	}
	if code == bimg.PDF {
		return "application/pdf"
	}
	if code == bimg.TIFF {
		return "image/tiff"
	}
	return "image/jpeg"
}
