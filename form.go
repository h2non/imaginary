package main

import "fmt"

func htmlForm() string {
	operations := []struct {
		name string
		args string
	}{
		{"resize", "width=300&height=200&type=png"},
		{"crop", "width=562&height=562&quality=95"},
		{"extract", "top=100&left=100&areawidth=300&areaheight=150"},
		{"enlarge", "width=1440&height=900&quality=95"},
		{"rotate", "rotate=180"},
		{"flip", ""},
		{"flop", ""},
		{"thumbnail", "width=100"},
		{"zoom", "factor=2&areawidth=300&top=80&left=80"},
		{"watermark", "textwidth=100&text=Hello&font=sans%2012&opacity=0.5&color=255,200,50"},
		{"convert", "type=png"},
		{"info", ""},
	}

	html := "<html><body>"

	for _, form := range operations {
		html += fmt.Sprintf(`
    <h1>%s</h1>
    <form method="POST" action="/%s?%s" enctype="multipart/form-data">
      <input type="file" name="file" />
      <input type="submit" value="Upload" />
    </form>`, form.name, form.name, form.args)
	}

	html += "</body></html>"

	return html
}
