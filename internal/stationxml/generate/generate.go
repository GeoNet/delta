package main

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

// A basic enumeration, will use uint32 as the base type.
var enumerationTemplate = `
{{- $base := . -}}
package stationxml

import (
	"encoding/xml"
	"fmt"
)

const (
	{{range $n, $e := .Enumerations -}}
	{{title (lower $e.Value)}}{{if eq 0 $n}} {{$base.AttrName}} = iota{{end}}
	{{end -}}
)

{{range $c := .Comments -}}
// {{$c}}
{{end -}}
type {{.AttrName}} uint32

func (v {{.AttrName}}) String() string {
	switch v {
        {{range $e := .Enumerations -}}
	case {{title (lower $e.Value)}}:
		return "{{$e.Value}}"
	{{end -}}
	default:
		return ""
	}
}

func (v {{.AttrName}}) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
        return xml.Attr{Name: name, Value: v.String()}, nil
}

func (v *{{.AttrName}}) UnmarshalXMLAttr(attr xml.Attr) error {
        switch attr.Value {
        {{range $e := .Enumerations -}}
        case "{{$e.Value}}":
                *v = {{title (lower $e.Value)}}
        {{end -}}
        default:
                return fmt.Errorf("unknown {{.AttrName}}: %s", attr.Value)
        }

        return nil
}

func ( v {{.AttrName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
        return e.EncodeElement(v.String(), start)
}

func (v *{{.AttrName}}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

        var s string
        err := d.DecodeElement(&s, &start)
        if err != nil {
                return err
        }

        switch s {
        {{range $e := .Enumerations -}}
        case "{{$e.Value}}":
                *v = {{title (lower $e.Value)}}
        {{end -}}
        default:
                return fmt.Errorf("unknown {{.AttrName}}: %s", s)
        }

        return nil
}
`

// A type that is an alias for a base type.
var simpleTemplate = `
{{- $base := . -}}
package stationxml

{{range $c := .Comments -}}
// {{$c}}
{{end -}}
type {{.AttrName}} {{type .Type .AttrName}}
`

// A type that is an extension of another type but only adds attribute variables.
var derivedTemplate = `
{{- $base := . -}}
package stationxml

{{range $c := .Comments -}}
// {{$c}}
{{end -}}
type {{.AttrName}} struct {
	{{type .Type .Type}}{{range $v := .Attributes }}

		{{range $c := $v.Comments -}}
			// {{$c}}
		{{end -}}
		{{title $v.AttrName}} {{type $v.AttrType "string"}} {{attr $v.AttrName $v.Use }}
	{{end -}}
}
`

// A type that can be an extension of another type but also may add attributes, or other variables.
var complexTemplate = `
{{- $base := . -}}
package stationxml

{{range $c := .Comments -}}
// {{$c}}
{{end -}}
type {{title .AttrName}} struct {
	{{with $v := .UseExtension }}

		{{type $v.Base $v.Base}}
	{{end -}}
	{{range $v := .Attributes }}

		{{range $c := $v.Comments -}}
			// {{$c}}
		{{end -}}
		{{title $v.AttrName}} {{type $v.AttrType "string"}} {{attr $v.AttrName $v.Use }}
	{{end -}}
	{{range $e := .Elements}}

		{{range $c := $e.Comments -}}
			// {{$c}}
		{{end -}}
		{{title $e.AttrName}} {{if $e.IsMultiple}}[]{{end}}{{type $e.Type $e.AttrName}} {{if $e.IsOptional}}{{omit $e.AttrName}}{{else}}{{xml $e.AttrName}}{{end}}
	{{end -}}
}
`

func (e *Element) Render(w io.Writer, tmpl string) error {

	t, err := template.New("element").Funcs(
		template.FuncMap{
			"title":  func(s string) string { return strings.Title(s) },
			"lower":  func(s string) string { return strings.ToLower(s) },
			"upper":  func(s string) string { return strings.ToUpper(s) },
			"suffix": func(s, v string) string { return strings.TrimSuffix(s, v) },
			"omit": func(s string, p ...string) string {
				switch {
				case len(p) > 0:
					return fmt.Sprintf("`xml:\"%s,%s,omitempty\"`", s, strings.Join(p, ","))
				default:
					return fmt.Sprintf("`xml:\"%s,omitempty\"`", s)
				}
			},
			"xml": func(s string, p ...string) string {
				switch {
				case len(p) > 0:
					return fmt.Sprintf("`xml:\"%s,%s\"`", s, strings.Join(p, ","))
				default:
					return fmt.Sprintf("`xml:\"%s\"`", s)
				}
			},
			"attr": func(s, u string) string {
				switch u {
				case "optional":
					return fmt.Sprintf("`xml:\"%s,attr,omitempty\"`", s)
				default:
					return fmt.Sprintf("`xml:\"%s,attr\"`", s)
				}
			},
			"type": func(s, t string) string {
				switch s {
				case "xs:integer":
					return "int"
				case "xs:string":
					return "string"
				case "xs:double", "xs:decimal":
					return "float64"
				case "xs:dateTime":
					return "DateTime"
				case "xs:anyURI":
					return "AnyURI"
				case "xs:NMTOKEN":
					return t
				}
				switch {
				case strings.HasPrefix(s, "fsx:"):
					return strings.TrimPrefix(s, "fsx:")
				default:
					return s
				}
			},
		},
	).Parse(tmpl)
	if err != nil {
		return err
	}
	if err := t.Execute(w, e); err != nil {
		return err
	}

	return nil
}
