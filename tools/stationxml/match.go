package main

import (
	"regexp"
	"strings"
)

type Matcher interface {
	MatchString(string) bool
}

type positiveMatch struct {
	re *regexp.Regexp
}

func (p positiveMatch) MatchString(s string) bool {
	if p.re != nil {
		return p.re.MatchString(s)
	}
	return true
}

type negativeMatch struct {
	re *regexp.Regexp
}

func (n negativeMatch) MatchString(s string) bool {
	if n.re != nil {
		return !n.re.MatchString(s)
	}
	return false
}

func Match(str string) (Matcher, error) {
	// prepare the matching string
	s := strings.TrimPrefix(strings.TrimSpace(str), "!")

	// compile the expression
	re, err := regexp.Compile(s)
	if err != nil {
		return nil, err
	}

	// a negative match if it starts with an "!" symbol
	if strings.HasPrefix(strings.TrimSpace(str), "!") {
		return negativeMatch{re}, nil
	}

	return positiveMatch{re}, nil
}

func MustMatch(str string) Matcher {
	m, err := Match(str)
	if err != nil {
		panic(err)
	}
	return m
}
