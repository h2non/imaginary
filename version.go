package main

import "gopkg.in/h2non/bimg.v1"

// Version stores the current package semantic version
const Version = "1.0.15"

// Versions represents the used versions for several significant dependencies
type Versions struct {
	ImaginaryVersion string `json:"imaginary"`
	BimgVersion      string `json:"bimg"`
	VipsVersion      string `json:"libvips"`
}

// CurrentVersions stores the current runtime system version metadata
var CurrentVersions = Versions{Version, bimg.Version, bimg.VipsVersion}
