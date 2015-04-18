package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const fixture = "fixtures/large.jpg"

func TestReadParams(t *testing.T) {
	var params ImageOptions

	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		params = readParams(r)
		w.Write([]byte{0})
	}

	url := "http://foo/?width=100&height=100&opacity=0.2&noreplicate=true&text=hello"
	file, _ := os.Open(fixture)
	r, _ := http.NewRequest("GET", url, file)
	w := httptest.NewRecorder()
	fakeHandler(w, r)

	assert := params.Width == 100 &&
		params.Height == 100 &&
		params.NoReplicate == true &&
		params.Opacity == 0.2 &&
		params.Text == "hello"

	if assert == false {
		t.Error("Invalid param")
	}
}

func TestParseColor(t *testing.T) {
	cases := []struct {
		value    string
		expected []uint8
	}{
		{"200,100,20", []uint8{200, 100, 20}},
		{"0,280,200", []uint8{0, 255, 200}},
		{" -1, 256 , 50", []uint8{0, 255, 50}},
		{" a, 20 , &hel0", []uint8{0, 20, 0}},
		{"", []uint8{}},
	}

	for _, color := range cases {
		c := parseColor(color.value)
		l := len(color.expected)

		if len(c) != l {
			t.Errorf("Invalid color length: %#v", c)
		}
		if l == 0 {
			continue
		}

		assert := c[0] == color.expected[0] &&
			c[1] == color.expected[1] &&
			c[2] == color.expected[2]

		if assert == false {
			t.Errorf("Invalid color schema: %#v <> %#v", color.expected, c)
		}
	}
}
