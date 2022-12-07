package main

import (
	"regexp"
)

// Matcher can be used as a cli argument.
type Matcher struct {
	*regexp.Regexp
}

// NewMatcher compiles the regexp and returns a Matcher.
func NewMatcher(exp string) (Matcher, error) {
	re, err := regexp.Compile(exp)
	if err != nil {
		return Matcher{}, err
	}
	return Matcher{re}, nil
}

// NewMatcher compiles the regexp and returns a Matcher, it will panic if an error is returned.
func MustMatcher(exp string) Matcher {
	re, err := NewMatcher(exp)
	if err != nil {
		panic(err)
	}
	return re
}

// MarshalText implements the TextMarshaller interface.
func (m Matcher) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// UnmarshalText implements the TextUnmarshaller interface.
func (m *Matcher) UnmarshalText(data []byte) error {
	re, err := regexp.Compile(string(data))
	if err != nil {
		return err
	}
	m.Regexp = re
	return nil
}
