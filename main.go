package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	Server(8088)
}

/*
func main() {
	options := vips.Options{
		Width:        300,
		Height:       240,
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      95,
	}

	f, _ := os.Open("fixtures/large.jpg")
	inBuf, _ := ioutil.ReadAll(f)

	buf, err := vips.Resize(inBuf, options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	err = ioutil.WriteFile("image.jpg", buf, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic("Cannot write file")
	}
}
*/
