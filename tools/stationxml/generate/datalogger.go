package main

import (
	"io"
	"strings"
	"text/template"
)

type DataloggerModel struct {
	Type         string `yaml:"type"`
	Description  string `yaml:"description"`
	Manufacturer string `yaml:"manufacturer"`
	Vendor       string `yaml:"vendor"`
}

var dataloggerModelTemplate = `

var DataloggerModels map[string]DataloggerModel = map[string]DataloggerModel{
{{ range $k, $v := . }}	"{{ $k}}": DataloggerModel{
		Type: "{{$v.Type}}",
		Description: "{{$v.Description}}",
		Manufacturer: "{{$v.Manufacturer}}",
		Vendor: "{{$v.Vendor}}",
	},
{{end}}{{"}"}}
`

func dataloggermodel(w io.Writer, dl map[string]DataloggerModel) error {

	t, err := template.New("dataloggermodels").Funcs(
		template.FuncMap{
			"escape": func(s string) string { return strings.Join(strings.Fields(s), " ") },
		},
	).Parse(dataloggerModelTemplate)
	if err != nil {
		return err
	}
	if err := t.Execute(w, dl); err != nil {
		return err
	}

	return nil
}
