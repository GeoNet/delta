package main

import (
	"aqwari.net/xml/xsd"
)

// ParseSelf decode the schema to process complex types
func (s *Schemas) ParseSelf(defaults Self, fn func(Self) error) error {

	for _, v := range s.Schema {
		if v.TargetNS != s.Namespace {
			continue
		}

		for k, t := range v.Types {
			if k.Space != s.Namespace {
				continue
			}

			switch t := t.(type) {
			case *xsd.ComplexType:
				if xsd.XMLName(t).Local != "_self" {
					continue
				}

				if err := fn(toSelf(t, defaults)); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

// ParseEnum decode the schema to process enum types
func (s *Schemas) ParseEnum(pkg string, fn func(string, Enum) error) error {

	for _, v := range s.Schema {
		if v.TargetNS != s.Namespace {
			continue
		}

		for k, t := range v.Types {
			if k.Space != s.Namespace {
				continue
			}

			switch t := t.(type) {
			case *xsd.SimpleType:
				if !(len(t.Restriction.Enum) > 0) {
					continue
				}

				if err := fn(xsd.XMLName(t).Local, toEnum(pkg, t)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// ParseSimple decode the schema to process simple types
func (s *Schemas) ParseSimple(pkg string, fn func(string, Simple) error) error {

	for _, v := range s.Schema {
		if v.TargetNS != s.Namespace {
			continue
		}

		for k, t := range v.Types {
			if k.Space != s.Namespace {
				continue
			}

			switch t := t.(type) {
			case *xsd.SimpleType:
				if len(t.Restriction.Enum) > 0 {
					continue
				}

				if err := fn(xsd.XMLName(t).Local, toSimple(pkg, t)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// ParseComplex decode the schema to process complex types
func (s *Schemas) ParseComplex(pkg string, fn func(string, Complex) error) error {

	for _, v := range s.Schema {
		if v.TargetNS != s.Namespace {
			continue
		}

		for k, t := range v.Types {
			if k.Space != s.Namespace {
				continue
			}

			switch t := t.(type) {
			case *xsd.ComplexType:
				if xsd.XMLName(t).Local == "_self" {
					continue
				}

				if err := fn(xsd.XMLName(t).Local, toComplex(pkg, t, s.Pointers...)); err != nil {
					return err
				}
			}

		}
	}

	return nil
}
