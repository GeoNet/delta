package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func main() {

	var dir string
	flag.StringVar(&dir, "dir", "responses", "response YAML directory")

	var xml string
	flag.StringVar(&xml, "xml", "auto", "output response XML directory")

	flag.Parse()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal(err)
	}

	if err := os.MkdirAll(xml, 0755); err != nil {
		log.Fatal(err)
	}

	response := NewResponseInfo()

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var r ResponseInfo
			if err := yaml.Unmarshal(b, &r); err != nil {
				return err
			}
			response.Merge(r)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	if err := response.Generate(os.Stdout); err != nil {
		log.Fatal(err)
	}

	if err := response.Build(xml); err != nil {
		log.Fatal(err)
	}
}
