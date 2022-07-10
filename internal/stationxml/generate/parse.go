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
