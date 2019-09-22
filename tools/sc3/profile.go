package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

type Profiler interface {
	Id() string
	Path() string
	Template() string
}

func Store(profile Profiler, path string) error {
	name := filepath.Join(path, profile.Path())

	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	tmpl, err := template.New("profile").Parse(profile.Template())
	if err != nil {
		return err
	}

	var res bytes.Buffer
	if err := tmpl.Execute(&res, profile); err != nil {
		return err
	}

	if err := ioutil.WriteFile(name, res.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
