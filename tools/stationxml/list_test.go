package main

import (
	"testing"
)

func TestListRegexp_File(t *testing.T) {

	var tests = []struct {
		f, r string
	}{
		{"testdata/list", "^(A1|A2|B1|B2|C1)$"},
	}

	for _, v := range tests {
		s, err := loadRegexpList(v.f)
		switch {
		case err != nil:
			t.Errorf("%s: unable to load regexp list file: %s", v.f, err)
		case s == nil:
			t.Errorf("%s: unable to load regexp list file: empty string", v.f)
		case string(s) != v.r:
			t.Errorf("%s: unable to load regexp list file: %s != %s", v.f, string(s), v.r)
		}
	}
}
