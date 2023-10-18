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

func ToKey(network, station string) string {
	return fmt.Sprintf("%s_%s", strings.ToUpper(network), strings.ToUpper(station))
}

type Station struct {
	Network string
	Code    string
}

func (s Station) Key() string {
	return ToKey(s.Network, s.Code)
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
