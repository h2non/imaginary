package main

import "fmt"

func htmlForm() string {
	operations := []struct {
		name   string
		method string
		args   string
	}{
		{"Resize", "resize", "width=300&height=200&type=png"},
		{"Force resize", "resize", "width=300&height=200&force=true"},
		{"Crop", "crop", "width=562&height=562&quality=95"},
		{"Extract", "extract", "top=100&left=100&areawidth=300&areaheight=150"},
		{"Enlarge", "enlarge", "width=1440&height=900&quality=95"},
		{"Rotate", "rotate", "rotate=180"},
		{"Flip", "flip", ""},
		{"Flop", "flop", ""},
		{"Thumbnail", "thumbnail", "width=100"},
		{"Zoom", "zoom", "factor=2&areawidth=300&top=80&left=80"},
		{"Color space (black&white)", "resize", "width=400&height=300&colorspace=bw"},
		{"Add watermark", "watermark", "textwidth=100&text=Hello&font=sans%2012&opacity=0.5&color=255,200,50"},
		{"Convert format", "convert", "type=png"},
		{"Image metadata", "info", ""},
	}

	html := "<html><body>"

	for _, form := range operations {
		html += fmt.Sprintf(`
    <h1>%s</h1>
    <form method="POST" action="/%s?%s" enctype="multipart/form-data">
      <input type="file" name="file" />
      <input type="submit" value="Upload" />
    </form>`, form.name, form.method, form.args)
	}

	html += "</body></html>"

	return html
}
