package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const contents = `###
### Delivered by puppet
###
capslink:sl4caps
`

type Station struct {
	Network string
	Code    string
}

func NewStation(s string) Station {
	r := strings.Split(strings.ToUpper(s), "_")
	switch n := len(r); {
	case n > 1:
		return Station{
			Network: r[0],
			Code:    r[1],
		}
	case n == 1:
		return Station{
			Network: "*",
			Code:    r[0],
		}
	default:
		return Station{
			Network: "*",
			Code:    "*",
		}
	}
}

func (s Station) Key() string {
	return fmt.Sprintf("%s_%s", strings.ToUpper(s.Network), strings.ToUpper(s.Code))
}

func (s Station) Path() string {
	return fmt.Sprintf("station_%s", s.Key())
}

func (s Station) Store(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		return err
	}
	return nil
}

func (s Station) Output(base string) error {
	if err := s.Store(filepath.Join(base, s.Path())); err != nil {
		return err
	}
	return nil
}
