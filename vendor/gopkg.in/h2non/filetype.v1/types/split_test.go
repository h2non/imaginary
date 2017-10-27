package types

import "testing"

func TestSplit(t *testing.T) {
	cases := []struct {
		mime    string
		kind    string
		subtype string
	}{
		{"image/jpeg", "image", "jpeg"},
		{"/jpeg", "", "jpeg"},
		{"image/", "image", ""},
		{"/", "", ""},
		{"image", "image", ""},
	}

	for _, test := range cases {
		kind, subtype := splitMime(test.mime)
		if test.kind != kind {
			t.Fatalf("Invalid kind: %s", test.kind)
		}
		if test.subtype != subtype {
			t.Fatalf("Invalid subtype: %s", test.subtype)
		}
	}
}
