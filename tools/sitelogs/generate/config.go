package main

import (
	"encoding/hex"
	"io"
	"text/template"
)

type Contact struct {
	Name               string `yaml:"Name"`
	TelephonePrimary   string `yaml:"TelephonePrimary"`
	TelephoneSecondary string `yaml:"TelephoneSecondary"`
	Fax                string `yaml:"Fax"`
	Email              string `yaml:"Email"`
}

type Agency struct {
	Agency                string  `yaml:"Agency"`
	PreferredAbbreviation string  `yaml:"PreferredAbbreviation"`
	MailingAddress        string  `yaml:"MailingAddress"`
	PrimaryContact        Contact `yaml:"PrimaryContact"`
	SecondaryContact      Contact `yaml:"SecondaryContact"`
	Notes                 string  `yaml:"Notes"`
}

type Country struct {
	Latitude  float64 `yaml:"Latitude"`
	Longitude float64 `yaml:"Longitude"`
}

type Config struct {
	PreparedBy            string             `yaml:"PreparedBy"`
	PrimaryDatacentre     string             `yaml:"PrimaryDatacentre"`
	URLForMoreInformation string             `yaml:"URLForMoreInformation"`
	ExtraNotes            string             `yaml:"ExtraNotes"`
	ContactAgency         Agency             `yaml:"ContactAgency"`
	ResponsibleAgency     Agency             `yaml:"ResponsibleAgency"`
	Countries             map[string]Country `yaml:"Countries"`
	Diagrams              map[string]string  `yaml:"Diagrams"`
}

var configTemplate = `
package main

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  To update: edit or add yaml file(s) in the config directory.
 *  Commit these changes and run "go generate" in the main project
 *  directory. Changes to this file should then also be commited.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

var preparedBy = "{{.PreparedBy}}"

var primaryDatacentre = {{tick}}{{.PrimaryDatacentre}}{{tick}}
var urlForMoreInformation = {{tick}}{{.URLForMoreInformation}}{{tick}}
var extraNotes = {{tick}}{{.ExtraNotes}}{{tick}}

{{define "contact"}}Contact{
Name:               {{tick}}{{.Name}}{{tick}},
TelephonePrimary:   {{tick}}{{.TelephonePrimary}}{{tick}},
TelephoneSecondary: {{tick}}{{.TelephoneSecondary}}{{tick}},
Fax:                {{tick}}{{.Fax}}{{tick}},
Email:              {{tick}}{{.Email}}{{tick}},
}{{end}}

{{define "agency"}}Agency{
Agency:                {{tick}}{{.Agency}}{{tick}},
PreferredAbbreviation: {{tick}}{{.PreferredAbbreviation}}{{tick}},
MailingAddress:        {{tick}}{{.MailingAddress}}{{tick}},
PrimaryContact:        {{template "contact" .PrimaryContact}},
SecondaryContact:      {{template "contact" .SecondaryContact}},
Notes:                 {{tick}}{{.Notes}}{{tick}},
}{{end}}

var contactAgency = {{template "agency" .ContactAgency}}

var responsibleAgency = {{template "agency" .ResponsibleAgency}}

var countryList = []struct {
	name     string
	lat, lon float64
}{ {{range $k, $v := .Countries}}
	{"{{$k}}", {{$v.Latitude}}, {{$v.Longitude}}{{"}"}},{{end}}
}

var antennaGraphs = map[string]string{
	{{ range $k, $v := .Diagrams }}    "{{$k}}": "{{hex $v}}",
{{end}}
}
`

func (c Config) Generate(w io.Writer) error {

	t, err := template.New("config").Funcs(
		template.FuncMap{
			"tick": func() string { return "`" },
			"hex":  func(s string) string { return hex.EncodeToString([]byte(s)) },
		},
	).Parse(configTemplate)
	if err != nil {
		return err
	}
	if err := t.Execute(w, c); err != nil {
		return err
	}

	return nil
}
