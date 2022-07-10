package main

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed tmpl/complex.tmpl
var complexTmpl string

var complexTemplate = template.Must(template.New("complex").Funcs(
	template.FuncMap{
		"bt": func() string { return "`" },
	}).Parse(complexTmpl))

type Variable struct {
	Name     string
	Type     string
	Pointer  bool
	Multiple bool
	Required bool
}

func (v Variable) Title() string {
	return strings.Title(v.Name)
}

func (v Variable) Optional() bool {
	return !v.Required
}

type Complex struct {
	Package string
	Name    string

	Derived *Variable
	Builtin *Variable

	Attributes []Variable
	Variables  []Variable
}

func (c Complex) Title() string {
	return strings.Title(c.Name)
}

func (c Complex) Render() ([]byte, error) {
	var buf bytes.Buffer
	if err := complexTemplate.Execute(&buf, c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c Complex) RenderFile(path string) error {
	return RenderFile(path, c)
}
