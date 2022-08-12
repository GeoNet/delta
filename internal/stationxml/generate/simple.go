package main

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed tmpl/simple.tmpl
var simpleTmpl string

var simpleTemplate = template.Must(template.New("simple").Funcs(
	template.FuncMap{
		"bt": func() string { return "`" },
	}).Parse(simpleTmpl))

type Simple struct {
	Package string
	Name    string
	Type    string
}

func (s Simple) Render() ([]byte, error) {
	var buf bytes.Buffer
	if err := simpleTemplate.Execute(&buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s Simple) RenderFile(path string) error {
	return RenderFile(path, s)
}
