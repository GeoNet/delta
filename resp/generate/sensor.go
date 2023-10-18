package main

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

type SensorComponent struct {
	Azimuth float64 `yaml:"azimuth"`
	Dip     float64 `yaml:"dip"`
}

type SensorModel struct {
	Type         string `yaml:"type"`
	Description  string `yaml:"description"`
	Manufacturer string `yaml:"manufacturer"`
	Vendor       string `yaml:"vendor"`

	Components []SensorComponent `yaml:"components"`
}

func (s SensorModel) Desc(name string) string {
	return fmt.Sprintf("%s %s %s", s.Manufacturer, strings.Split(strings.Fields(name)[0], "/")[0], s.Type)
}

func (s SensorModel) Make() string {
	switch parts := strings.Fields(s.Manufacturer); {
	case len(parts) > 0:
		label := strings.Join(parts, " ")
		for _, s := range []string{"/", ".", "+", " "} {
			label = strings.ReplaceAll(label, s, "-")
		}
		return strings.TrimRight(label, "-")
	default:
		return "Unknown"
	}
}

var sensorModelTemplate = `

var SensorModels map[string]SensorModel = map[string]SensorModel{
{{ range $k, $v := . }}	"{{ $k}}": SensorModel{
		Name: "{{$k}}",
		Type: "{{$v.Type}}",
		Description: "{{$v.Desc $k}}",
		Manufacturer: "{{$v.Manufacturer}}",
		Vendor: "{{$v.Vendor}}",
		Components: []SensorComponent{{"{"}}{{ range $z := $v.Components}}SensorComponent{Azimuth: {{ $z.Azimuth }}, Dip: {{ $z.Dip }}{{"}"}},{{end}}{{"}"}},
	},
{{end}}{{"}"}}
`

func sensormodel(w io.Writer, dl map[string]SensorModel) error {

	t, err := template.New("sensormodels").Funcs(
		template.FuncMap{
			"escape": func(s string) string { return strings.Join(strings.Fields(s), " ") },
		},
	).Parse(sensorModelTemplate)
	if err != nil {
		return err
	}
	if err := t.Execute(w, dl); err != nil {
		return err
	}

	return nil
}
