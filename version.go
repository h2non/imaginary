package main

import "gopkg.in/h2non/bimg.v1"

<<<<<<< HEAD
const Version = "0.1.25"
=======
const Version = "0.1.24"
>>>>>>> 0087d8d343dc24b81b66ac0932030310267c1aa6

type Versions struct {
	ImaginaryVersion string `json:"imaginary"`
	BimgVersion      string `json:"bimg"`
	VipsVersion      string `json:"libvips"`
}

var CurrentVersions = Versions{Version, bimg.Version, bimg.VipsVersion}
