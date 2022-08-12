package main

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed tmpl/self.tmpl
var selfTmpl string

var selfTemplate = template.Must(template.New("self").Funcs(
	template.FuncMap{
		"bt": func() string { return "`" },
	}).Parse(selfTmpl))

type Self struct {
	Package string
	Space   string
	Name    string
	Derived string
	Version float64
}

func (s Self) Render() ([]byte, error) {
	var buf bytes.Buffer
	if err := selfTemplate.Execute(&buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s Self) RenderFile(path string) error {
	return RenderFile(path, s)
}
