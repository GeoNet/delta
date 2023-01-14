package main

import (
	"regexp"
)

type Matcher struct {
	*regexp.Regexp
}

func NewMatcher(exp string) (Matcher, error) {
	re, err := regexp.Compile(exp)
	if err != nil {
		return Matcher{}, err
	}
	return Matcher{re}, nil
}

func MustMatcher(exp string) Matcher {
	re, err := NewMatcher(exp)
	if err != nil {
		panic(err)
	}
	return re
}

func (m Matcher) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Matcher) UnmarshalText(data []byte) error {
	re, err := regexp.Compile(string(data))
	if err != nil {
		return err
	}
	m.Regexp = re
	return nil
}
