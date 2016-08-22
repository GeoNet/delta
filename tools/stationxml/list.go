package main

import (
	"io/ioutil"
	"strings"
)

func loadRegexpList(path string) ([]byte, error) {

	s, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return []byte("^(" + strings.Join(strings.Fields(string(s)), "|") + ")$"), nil
}
