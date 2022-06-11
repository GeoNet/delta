package main

import (
	"io"
	"strings"
	"text/template"
)

var generateTemplate = `
/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  To update any changes, run "go generate" in the main project
 *  directory and then commit this file.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

package meta

{{range $k, $v := .Fields -}}
// {{title $k}} is a helper function to return a slice copy of {{$v.Key}} values.
func (s Set) {{title $k}}() []{{$v.Key}} {
	{{$k}} := make([]{{$v.Key}}, len(s.{{$k}}))
        copy({{$k}}, s.{{$k}})
        return {{$k}}
}

{{end -}}

{{range $k, $v := .Lookup -}}
// {{$v.Key}} is a helper function to return a {{$v.Key}} value and true if one exists.
func (s Set) {{$v.Key}}({{join ", " $v.Fields}} string) ({{$v.Key}}, bool) {
	for _, v := range s.{{$k}} {
		{{range $f := $v.Fields -}}
		if {{$f}} != v.{{title $f}} {
			continue
		}
		{{end -}}
		return v, true
	}
	return {{$v.Key}}{{"{}"}}, false
}

{{end -}}
`

type Generate struct {
	Fields map[string]struct {
		Key string //nolint:unused // used in template
	}
	Lookup map[string]struct {
		Key    string   //nolint:unused // used in template
		Fields []string //nolint:unused // used in template
	}
}

func (g Generate) Write(w io.Writer) error {

	t, err := template.New("generate").Funcs(
		template.FuncMap{
			"title": func(s string) string { return strings.Title(s) },
			"join":  func(k string, s []string) string { return strings.Join(s, k) },
		},
	).Parse(generateTemplate)
	if err != nil {
		return err
	}
	if err := t.Execute(w, g); err != nil {
		return err
	}

	return nil
}
