package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func main() {

	var dir string
	flag.StringVar(&dir, "dir", "responses", "response YAML directory")

	flag.Parse()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal(err)
	}

	response := NewResponseInfo()

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			b, err := ioutil.ReadFile(path)
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
}
