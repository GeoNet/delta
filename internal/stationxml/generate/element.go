package main

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

type Element struct {
	AttrName        string `xml:"name,attr,omitempty"`
	AttrType        string `xml:"type,attr,omitempty"`
	Use             string `xml:"use,attr,omitempty"`
	Default         string `xml:"default,attr,omitempty"`
	Fixed           string `xml:"fixed,attr,omitempty"`
	Ref             string `xml:"ref,attr,omitempty"`
	Base            string `xml:"base,attr,omitempty"`
	Value           string `xml:"value,attr,omitempty"`
	NameSpace       string `xml:"namespace,attr,omitempty"`
	ProcessContents string `xml:"processContents,attr,omitempty"`

	MinOccurs    string   `xml:"minOccurs,attr,omitempty"`
	MaxOccurs    string   `xml:"maxOccurs,attr,omitempty"`
	MinInclusive *Element `xml:"minInclusive,omitempty"`
	MaxInclusive *Element `xml:"maxInclusive,omitempty"`
	MaxExclusive *Element `xml:"maxExclusive,omitempty"`

	Annotation *Annotation `xml:"annotation,omitempty"`

	Restriction    *Element `xml:"restriction,omitempty"`
	SimpleType     *Element `xml:"simpleType,omitempty"`
	ComplexType    *Element `xml:"complexType,omitempty"`
	Pattern        *Element `xml:"patter,omitempty"`
	SimpleContent  *Element `xml:"simpleContent,omitempty"`
	ComplexContent *Element `xml:"complexContent,omitempty"`
	Extension      *Element `xml:"extension,omitempty"`

	Attribute      []*Element `xml:"attribute,omitempty"`
	AttributeGroup []*Element `xml:"attributeGroup,omitempty"`
	Enumeration    []*Element `xml:"enumeration,omitempty"`
	Sequence       []*Element `xml:"sequence,omitempty"`
	Group          []*Element `xml:"group,omitempty"`
	Choice         []*Element `xml:"choice,omitempty"`
	Element        []*Element `xml:"element,omitempty"`
	AnyAttribute   []*Element `xml:"anyAttribute,omitempty"`
}

func (e *Element) Type() string {

	if r := e.Restriction; r != nil && r.Base != "" {
		return r.Base
	}

	if t := e.SimpleType; t != nil {
		if v := t.Restriction; v != nil && v.Base != "" {
			return v.Base
		}
	}

	if t := e.ComplexType; t != nil {
		if s := t.SimpleContent; s != nil {
			if v := s.Extension; v != nil && v.Base != "" {
				return v.Base
			}
			if v := s.Restriction; v != nil && v.Base != "" {
				return v.Base
			}
		}
	}

	if t := e.AttrType; t != "" {
		return t
	}

	return e.AttrName
}

func (e *Element) Comments() []string {
	if a := e.Annotation; a != nil {
		return a.Comments()
	}
	return nil
}

func (e *Element) IsEnumeration() bool {
	return e.Type() == "xs:NMTOKEN"
}

func (e *Element) IsSimple() bool {
	switch t := e.Type(); {
	case t == "xs:NMTOKEN":
		return false
	case strings.HasPrefix(t, "xs:"):
		return true
	default:
		return false
	}
}

func (e *Element) IsDerived() bool {
	return strings.HasPrefix(e.Type(), "fsx:")
}

func (e *Element) UseExtension() *Element {
	if c := e.ComplexContent; c != nil {
		if t := c.Extension; t != nil && t.Base != "" {
			return t
		}
	}
	return nil
}

func (e *Element) Enumerations() []*Element {
	var elements []*Element
	e.SearchAll(func(key string, self *Element) {
		if key != "enumeration" {
			return
		}
		elements = append(elements, self)
	})
	return elements
}

func (e *Element) Attributes() []*Element {
	var elements []*Element
	e.Search(func(key string, self *Element) bool {
		if key == "element" && self.AttrType == "" {
			return true
		}
		if key != "attribute" {
			return false
		}
		elements = append(elements, self)
		return false
	})
	return elements
}

func (e *Element) IsMultiple() bool {
	return e.MaxOccurs == "unbounded"
}

func (e *Element) IsOptional() bool {
	return e.MinOccurs == "0"
}

func (e *Element) Elements() []*Element {
	var elements []*Element
	e.Search(func(key string, self *Element) bool {
		if key != "element" {
			return false
		}
		elements = append(elements, self)
		if self.AttrType == "" {
			return true
		}
		return false
	})
	return elements
}

func (e *Element) Search(fn func(string, *Element) bool) {
	e.Walk("", fn)
}

func (e *Element) SearchAll(fn func(string, *Element)) {
	e.WalkAll("", fn)
}

func (e *Element) WalkAll(key string, fn func(string, *Element)) {
	e.Walk(key, func(key string, self *Element) bool {
		fn(key, self)
		return false
	})
}

func (e *Element) Walk(key string, fn func(string, *Element) bool) {

	// process the function,
	// stop recursing if true is returned
	if fn(key, e) {
		return
	}

	if v := e.MinInclusive; v != nil {
		v.Walk("minInclusive", fn)
	}

	if v := e.MaxInclusive; v != nil {
		v.Walk("maxInclusive", fn)
	}

	if v := e.MaxExclusive; v != nil {
		v.Walk("maxExclusive", fn)
	}

	if v := e.Restriction; v != nil {
		v.Walk("restriction", fn)
	}

	if v := e.SimpleType; v != nil {
		v.Walk("simpleType", fn)
	}

	if v := e.ComplexType; v != nil {
		v.Walk("complexType", fn)
	}

	if v := e.Pattern; v != nil {
		v.Walk("pattern", fn)
	}

	if v := e.SimpleContent; v != nil {
		v.Walk("simpleContent", fn)
	}

	if v := e.ComplexContent; v != nil {
		v.Walk("complexContent", fn)
	}

	if v := e.Extension; v != nil {
		v.Walk("extension", fn)
	}

	for _, v := range e.Attribute {
		v.Walk("attribute", fn)
	}

	for _, v := range e.AttributeGroup {
		v.Walk("attributeGroup", fn)
	}

	for _, v := range e.Enumeration {
		v.Walk("enumeration", fn)
	}

	for _, v := range e.Sequence {
		v.Walk("sequence", fn)
	}

	for _, v := range e.Group {
		v.Walk("group", fn)
	}

	for _, v := range e.Choice {
		v.Walk("choice", fn)
	}

	for _, v := range e.Element {
		v.Walk("element", fn)
	}

	for _, v := range e.AnyAttribute {
		v.Walk("anyAttribute", fn)
	}
}

func (e *Element) Render(fsys fs.FS, w io.Writer, tmpl string) error {
	t, err := template.New(filepath.Base(tmpl)).Funcs(
		template.FuncMap{
			"title":  func(s string) string { return strings.Title(s) },
			"lower":  func(s string) string { return strings.ToLower(s) },
			"upper":  func(s string) string { return strings.ToUpper(s) },
			"suffix": func(s, v string) string { return strings.TrimSuffix(s, v) },
			"use": func(s string) string {
				switch s {
				case "optional":
					return "omitempty"
				default:
					return ""
				}
			},
			"xml": func(s string, p ...string) string {
				list := []string{s}
				for _, v := range p {
					if x := strings.TrimSpace(v); x != "" {
						list = append(list, x)
					}
				}
				return fmt.Sprintf("`xml:\"%s\"`", strings.Join(list, ","))
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
	).ParseFS(fsys, tmpl)
	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(w, filepath.Base(tmpl), e); err != nil {
		return err
	}

	return nil
}
