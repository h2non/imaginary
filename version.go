package main

// Version stores the current package semantic version
var Version = "1.1.1"

// Versions represents the used versions for several significant dependencies
type Versions struct {
	ImaginaryVersion string `json:"imaginary"`
	BimgVersion      string `json:"bimg"`
	VipsVersion      string `json:"libvips"`
}
