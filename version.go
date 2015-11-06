package main

import "gopkg.in/h2non/bimg.v0"

const Version = "0.1.18"

type Versions struct {
	ImaginaryVersion string `json:"imaginary"`
	BimgVersion      string `json:"bimg"`
	VipsVersion      string `json:"libvips"`
}

var CurrentVersions = Versions{Version, bimg.Version, bimg.VipsVersion}
