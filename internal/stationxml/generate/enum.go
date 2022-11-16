package main

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed tmpl/enum.tmpl
var enumTmpl string

var enumTemplate = template.Must(template.New("enum").Funcs(
	template.FuncMap{
		"bt": func() string { return "`" },
	}).Parse(enumTmpl))

type Enum struct {
	Package string
	Type    string
	Values  []string
}

func (e Enum) Name(v string) string {
	parts := strings.FieldsFunc(v, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	var list []string
	for _, p := range parts {
		list = append(list, cases.Title(language.English).String(strings.ToLower(p)))
	}
	return strings.TrimSuffix(strings.Join(list, "")+e.Type, "Type")
}

func (e Enum) Render() ([]byte, error) {
	var buf bytes.Buffer
	if err := enumTemplate.Execute(&buf, e); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e Enum) RenderFile(path string) error {
	return RenderFile(path, e)
}
