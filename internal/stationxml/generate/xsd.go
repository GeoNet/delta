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
