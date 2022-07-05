package main

import (
	"encoding/xml"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var name string
	flag.StringVar(&name, "name", "stationxml", "package name")

	var input string
	flag.StringVar(&input, "schema", "fdsn-station-1.2.xsd", "input schema file")

	var output string
	flag.StringVar(&output, "output", "output", "output dir")

	flag.Parse()

	if err := os.MkdirAll(output, 0755); err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}

	var schema Schema
	if err := xml.Unmarshal(data, &schema); err != nil {
		log.Fatal(err)
	}

	var elements []*Element

	elements = append(elements, schema.Groups()...)
	elements = append(elements, schema.AttributeGroups()...)
	elements = append(elements, schema.Simple()...)
	elements = append(elements, schema.Complex()...)
	elements = append(elements, schema.Elements()...)

	for _, e := range elements {
		switch path := filepath.Join(output, strings.Title(e.AttrName)); {
		case e.IsEnumeration():
			if err := RenderFile(path, enumerationTemplate, e); err != nil {
				log.Fatal(err)
			}
		case e.IsSimple():
			if err := RenderFile(path, simpleTemplate, e); err != nil {
				log.Fatal(err)
			}
		case e.IsDerived():
			if err := RenderFile(path, derivedTemplate, e); err != nil {
				log.Fatal(err)
			}
		default:
			if err := RenderFile(path, complexTemplate, e); err != nil {
				log.Fatal(e.AttrName, err)
			}
		}
	}
}
