package main

import (
	"io"
	"io/fs"
	"path/filepath"
	"text/template"
)

type Schema struct {
	TargetNamespace      string `xml:"targetNamespace,attr"`
	ElementFormDefault   string `xml:"elementFormDefault,attr"`
	AttributeFormDefault string `xml:"attributeFormDefault,attr"`
	Version              string `xml:"version,attr"`

	Annotation *Annotation `xml:"annotation,omitempty"`
	Element    *Element    `xml:"element,omitempty"`

	AttributeGroup []*Element `xml:"attributeGroup,omitempty"`
	SimpleType     []*Element `xml:"simpleType,omitempty"`
	ComplexType    []*Element `xml:"complexType,omitempty"`
	Group          []*Element `xml:"group,omitempty"`
}

func (s Schema) Comments() []string {
	if a := s.Annotation; a != nil {
		return a.Comments()
	}
	return nil
}

func (s Schema) Walk(fn func(string, *Element)) {
	if e := s.Element; e != nil {
		e.WalkAll("element", fn)
	}
	for _, v := range s.AttributeGroup {
		v.WalkAll("attributeGroup", fn)
	}
	for _, v := range s.SimpleType {
		v.WalkAll("simpleType", fn)
	}
	for _, v := range s.ComplexType {
		v.WalkAll("complexType", fn)
	}
	for _, v := range s.Group {
		v.WalkAll("group", fn)
	}
}

func (s Schema) Groups() []*Element {
	var elements []*Element
	s.Walk(func(element string, self *Element) {
		if element != "group" || self.AttrName == "" {
			return
		}
		elements = append(elements, self)
	})
	return elements
}

func (s Schema) AttributeGroups() []*Element {
	var elements []*Element
	s.Walk(func(element string, self *Element) {
		if element != "attributeGroup" || self.AttrName == "" {
			return
		}
		elements = append(elements, self)
	})
	return elements
}

func (s Schema) Simple() []*Element {
	var elements []*Element
	s.Walk(func(element string, self *Element) {
		if element != "simpleType" || self.AttrName == "" {
			return
		}
		elements = append(elements, self)
	})
	return elements
}

func (s Schema) Complex() []*Element {
	var elements []*Element
	s.Walk(func(element string, self *Element) {
		if element != "complexType" || self.AttrName == "" {
			return
		}
		elements = append(elements, self)
	})
	return elements
}

func (s Schema) Elements() []*Element {
	var elements []*Element
	s.Walk(func(element string, self *Element) {
		if element != "element" || self.AttrName == "" || self.AttrType != "" {
			return
		}
		elements = append(elements, self)
	})
	return elements
}

func (s Schema) Render(fsys fs.FS, w io.Writer, tmpl string) error {
	t, err := template.New(filepath.Base(tmpl)).Funcs(
		template.FuncMap{
			"bt": func() string { return "`" },
		},
	).ParseFS(fsys, tmpl)
	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(w, filepath.Base(tmpl), s); err != nil {
		return err
	}

	return nil
}
