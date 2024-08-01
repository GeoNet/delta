package main

import (
	"aqwari.net/xml/xsd"
)

func toFormat(t xsd.Type) string {
	switch t := t.(type) {
	case xsd.Builtin:
		switch t.String() {
		case "String", "NMTOKEN":
			return "string"
		case "Double":
			return "float64"
		case "DateTime":
			return "DateTime"
		case "Integer":
			return "int"
		case "Decimal":
			return "float64"
		default:
			return ""
		}
	default:
		return ""
	}
}

func toType(element xsd.Element) string {
	switch t := element.Type.(type) {
	case *xsd.SimpleType:
		return t.Name.Local
	case *xsd.ComplexType:
		return t.Name.Local
	default:
		return toFormat(element.Type)
	}
}

func toSelf(cmplx *xsd.ComplexType, defaults Self) Self {

	for _, e := range cmplx.Elements {
		return Self{
			Package: defaults.Package,
			Space:   e.Name.Space,
			Name:    e.Name.Local,
			Derived: toType(e),
			Version: defaults.Version,
		}
	}

	return defaults
}

func toEnum(pkg string, simple *xsd.SimpleType) Enum {
	return Enum{
		Package: pkg,
		Type:    simple.Name.Local,
		Values:  simple.Restriction.Enum,
	}
}

func toSimple(pkg string, simple *xsd.SimpleType) Simple {
	return Simple{
		Package: pkg,
		Name:    simple.Name.Local,
		Type:    toFormat(simple.Base),
	}
}

func toAttribute(attr xsd.Attribute) Variable {

	switch t := attr.Type.(type) {
	case *xsd.SimpleType:
		return Variable{
			Name:     attr.Name.Local,
			Type:     t.Name.Local,
			Required: !attr.Optional,
		}
	default:
		return Variable{
			Name:     attr.Name.Local,
			Type:     toFormat(attr.Type),
			Required: !attr.Optional,
		}
	}
}

func isPointer(element xsd.Element, pointers ...string) bool {

	for _, p := range pointers {
		if element.Name.Local != p {
			continue
		}
		return true
	}

	switch element.Type.(type) {
	case *xsd.SimpleType:
		return element.Optional && !element.Plural
	case *xsd.ComplexType:
		return element.Optional && !element.Plural
	default:
		return false
	}
}

func toVariable(element xsd.Element, pointers ...string) Variable {

	switch t := element.Type.(type) {
	case *xsd.SimpleType:
		return Variable{
			Name:     element.Name.Local,
			Type:     t.Name.Local,
			Required: !element.Optional,
			Multiple: element.Plural,
			Pointer:  isPointer(element, pointers...),
		}
	case *xsd.ComplexType:
		return Variable{
			Name:     element.Name.Local,
			Type:     t.Name.Local,
			Required: !element.Optional,
			Multiple: element.Plural,
			Pointer:  isPointer(element, pointers...),
		}
	default:
		return Variable{
			Name:     element.Name.Local,
			Type:     toFormat(element.Type),
			Required: !element.Optional,
			Pointer:  isPointer(element, pointers...),
			Multiple: element.Plural,
		}
	}
}

func toComplex(pkg string, cmplx *xsd.ComplexType, pointers ...string) Complex {

	var derived *Variable
	var builtin *Variable

	switch t := cmplx.Base.(type) {
	case *xsd.ComplexType:
		derived = &Variable{
			Name: t.Name.Local,
		}
	case xsd.Builtin:
		if v := toFormat(t); v != "" {
			builtin = &Variable{
				Type: v,
			}
		}
	}

	var attributes []Variable
	for _, a := range cmplx.Attributes {
		if v := toAttribute(a); v.Type != "" {
			attributes = append(attributes, v)
		}
	}

	var variables []Variable
	for _, e := range cmplx.Elements {
		if v := toVariable(e, pointers...); v.Type != "" {
			variables = append(variables, v)
		}
	}

	return Complex{
		Package:    pkg,
		Name:       cmplx.Name.Local,
		Derived:    derived,
		Builtin:    builtin,
		Attributes: attributes,
		Variables:  variables,
	}
}
