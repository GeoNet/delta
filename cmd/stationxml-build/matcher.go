package main

import (
	"regexp"
	"strings"
)

// Matcher can be used as a cli argument.
type Matcher struct {
	regexp *regexp.Regexp
	invert bool
}

// NewMatcher compiles the regexp and returns a Matcher.
func NewMatcher(exp string) (Matcher, error) {
	switch {
	case strings.HasPrefix(exp, "!"):
		re, err := regexp.Compile(exp[1:])
		if err != nil {
			return Matcher{}, err
		}
		return Matcher{regexp: re, invert: true}, nil
	default:
		re, err := regexp.Compile(exp)
		if err != nil {
			return Matcher{}, err
		}
		return Matcher{regexp: re, invert: false}, nil
	}
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
	switch {
	case m.invert:
		return []byte("!" + m.regexp.String()), nil
	default:
		return []byte(m.regexp.String()), nil
	}
}

// UnmarshalText implements the TextUnmarshaller interface.
func (m *Matcher) UnmarshalText(data []byte) error {
	switch s := string(data); {
	case strings.HasPrefix(s, "!"):
		re, err := regexp.Compile(s[1:])
		if err != nil {
			return err
		}
		m.regexp = re
		m.invert = true
	default:
		re, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		m.regexp = re
		m.invert = false
	}
	return nil
}

func (m Matcher) Match(check []byte) bool {
	switch ok := m.regexp.Match(check); {
	case m.invert:
		return !ok
	default:
		return ok
	}
}

func (m Matcher) MatchString(check string) bool {
	switch ok := m.regexp.MatchString(check); {
	case m.invert:
		return !ok
	default:
		return ok
	}
}
