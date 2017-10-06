package main

import (
	"net/url"
	"testing"

	"gopkg.in/h2non/bimg.v1"
)

const fixture = "fixtures/large.jpg"

func TestReadParams(t *testing.T) {
	q := url.Values{}
	q.Set("width", "100")
	q.Add("height", "80")
	q.Add("noreplicate", "1")
	q.Add("opacity", "0.2")
	q.Add("text", "hello")
	q.Add("background", "255,10,20")

	params := readParams(q)

	assert := params.Width == 100 &&
		params.Height == 80 &&
		params.NoReplicate == true &&
		params.Opacity == 0.2 &&
		params.Text == "hello" &&
		params.Background[0] == 255 &&
		params.Background[1] == 10 &&
		params.Background[2] == 20

	if assert == false {
		t.Error("Invalid params")
	}
}

func TestParseParam(t *testing.T) {
	intCases := []struct {
		value    string
		expected int
	}{
		{"1", 1},
		{"0100", 100},
		{"-100", 100},
		{"99.02", 99},
		{"99.9", 100},
	}

	for _, test := range intCases {
		val := parseParam(test.value, "int")
		if val != test.expected {
			t.Errorf("Invalid param: %s != %d", test.value, test.expected)
		}
	}

	floatCases := []struct {
		value    string
		expected float64
	}{
		{"1.1", 1.1},
		{"01.1", 1.1},
		{"-1.10", 1.10},
		{"99.999999", 99.999999},
	}

	for _, test := range floatCases {
		val := parseParam(test.value, "float")
		if val != test.expected {
			t.Errorf("Invalid param: %#v != %#v", val, test.expected)
		}
	}

	boolCases := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1", true},
		{"1.1", false},
		{"-1", false},
		{"0", false},
		{"0.0", false},
		{"no", false},
		{"yes", false},
	}

	for _, test := range boolCases {
		val := parseParam(test.value, "bool")
		if val != test.expected {
			t.Errorf("Invalid param: %#v != %#v", val, test.expected)
		}
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

func TestParseExtend(t *testing.T) {
	cases := []struct {
		value    string
		expected bimg.Extend
	}{
		{"white", bimg.ExtendWhite},
		{"black", bimg.ExtendBlack},
		{"copy", bimg.ExtendCopy},
		{"mirror", bimg.ExtendMirror},
		{"background", bimg.ExtendBackground},
		{" BACKGROUND  ", bimg.ExtendBackground},
		{"invalid", bimg.ExtendBlack},
		{"", bimg.ExtendBlack},
	}

	for _, extend := range cases {
		c := parseExtendMode(extend.value)
		if c != extend.expected {
			t.Errorf("Invalid extend value : %d != %d", c, extend.expected)
		}
	}
}

func TestGravity(t *testing.T) {
	cases := []struct {
		gravityValue   string
		smartCropValue bool
	}{
		{gravityValue: "foo", smartCropValue: false},
		{gravityValue: "smart", smartCropValue: true},
	}

	for _, td := range cases {
		io := readParams(url.Values{"gravity": []string{td.gravityValue}})
		if (io.Gravity == bimg.GravitySmart) != td.smartCropValue {
			t.Errorf("Expected %t to be %t, test data: %+v", io.Gravity == bimg.GravitySmart, td.smartCropValue, td)
		}
	}
}

func TestReadMapParams(t *testing.T) {
	cases := []struct {
		params   map[string]interface{}
		expected ImageOptions
	}{
		{
			map[string]interface{}{
				"width":   100,
				"opacity": 0.1,
				"type":    "webp",
				"embed":   true,
				"gravity": "west",
				"color":   "255,200,150",
			},
			ImageOptions{
				Width:   100,
				Opacity: 0.1,
				Type:    "webp",
				Embed:   true,
				Gravity: bimg.GravityWest,
				Color:   []uint8{255, 200, 150},
			},
		},
	}

	for _, test := range cases {
		opts := readMapParams(test.params)
		if opts.Width != test.expected.Width {
			t.Errorf("Invalid width: %d != %d", opts.Width, test.expected.Width)
		}
		if opts.Opacity != test.expected.Opacity {
			t.Errorf("Invalid opacity: %#v != %#v", opts.Opacity, test.expected.Opacity)
		}
		if opts.Type != test.expected.Type {
			t.Errorf("Invalid type: %s != %s", opts.Type, test.expected.Type)
		}
		if opts.Embed != test.expected.Embed {
			t.Errorf("Invalid embed: %#v != %#v", opts.Embed, test.expected.Embed)
		}
		if opts.Gravity != test.expected.Gravity {
			t.Errorf("Invalid gravity: %#v != %#v", opts.Gravity, test.expected.Gravity)
		}
		if opts.Color[0] != test.expected.Color[0] || opts.Color[1] != test.expected.Color[1] || opts.Color[2] != test.expected.Color[2] {
			t.Errorf("Invalid color: %#v != %#v", opts.Color, test.expected.Color)
		}
	}
}
