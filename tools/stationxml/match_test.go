package main

import (
	"testing"
)

func TestMatcher(t *testing.T) {

	var tests = map[string]struct {
		m Matcher
		s string
		r bool
	}{
		"empty positive":         {positiveMatch{}, "", true},
		"empty negative":         {negativeMatch{}, "", false},
		"valid positive match":   {MustMatch("AB"), "AB", true},
		"invalid positive match": {MustMatch("AB"), "BB", false},
		"valid negative match":   {MustMatch("!AB"), "BB", true},
		"invalid negative match": {MustMatch("!AB"), "AB", false},
	}

	for k, v := range tests {
		t.Log("test matcher for " + k)
		if v.m.MatchString(v.s) != v.r {
			t.Errorf("%s: unable to match string: %s", k, v.s)
		}
	}
}
