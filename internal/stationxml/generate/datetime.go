package main

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed tmpl/datetime.tmpl
var datetimeTmpl string

var datetimeTemplate = template.Must(template.New("datetime").Funcs(
	template.FuncMap{
		"bt": func() string { return "`" },
	}).Parse(datetimeTmpl))

type Datetime struct {
	Package string
	Format  string
	Future  bool
}

func (s Datetime) Render() ([]byte, error) {
	var buf bytes.Buffer
	if err := datetimeTemplate.Execute(&buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s Datetime) RenderFile(path string) error {
	return RenderFile(path, s)
}
