package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// download the remote url contents, an insecure option is available for problem sites.
func download(url string, insecure bool) ([]byte, error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, //nolint:gosec // needed until fdsn.org passed checks
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func main() {

	var name string
	flag.StringVar(&name, "name", "stationxml", "package name")

	var base string
	flag.StringVar(&base, "base", "base.go", "file name for extra go code")

	var doc string
	flag.StringVar(&doc, "doc", "doc.go", "file name for schema doc go code")

	var datetime string
	flag.StringVar(&datetime, "datetime", "date_time.go", "file name for extra date time go code")

	var test string
	flag.StringVar(&test, "test", "base_test.go", "file name for extra go test code")

	var version string
	flag.StringVar(&version, "version", "v1.2", "optional schema version to add to directory path")

	var input string
	flag.StringVar(&input, "input", "", "input schema file")

	var remote string
	flag.StringVar(&remote, "schema", "", "schema service endpoint to download from")

	var output string
	flag.StringVar(&output, "output", "output", "output dir")

	var insecure bool
	flag.BoolVar(&insecure, "insecure", false, "whether the remote site has certificate issues, use with caution")

	flag.Parse()

	var schema Schema
	switch {
	case remote != "":
		data, err := download(remote, insecure)
		if err != nil {
			log.Fatalf("unable to download schema from %s: %v", remote, err)
		}
		if err := xml.Unmarshal(data, &schema); err != nil {
			log.Fatalf("unable to unmarshal schema from %s: %v", remote, err)
		}
	case input != "":
		data, err := os.ReadFile(input)
		if err != nil {
			log.Fatalf("unable to read schema from input file %s: %v", remote, err)
		}
		if err := xml.Unmarshal(data, &schema); err != nil {
			log.Fatalf("unable to unmarshal schema from input file %s: %v", remote, err)
		}
	default:
		log.Fatal("no schema source found, needs either an input file or remote url")
	}

	if err := os.MkdirAll(output, 0755); err != nil {
		log.Fatal(err)
	}

	for k, v := range map[string]string{
		base:     baseTemplate,
		doc:      docTemplate,
		datetime: dateTimeTemplate,
		test:     testTemplate,
	} {
		if err := FormatFile(filepath.Join(output, k), v, schema); err != nil {
			log.Fatalf("unable to format %s: %v", k, err)
		}
	}

	var elements []*Element

	elements = append(elements, schema.Groups()...)
	elements = append(elements, schema.AttributeGroups()...)
	elements = append(elements, schema.Simple()...)
	elements = append(elements, schema.Complex()...)
	elements = append(elements, schema.Elements()...)

	for _, e := range elements {
		switch path := filepath.Join(output, FileName(strings.Title(e.AttrName), ".go")); {
		case e.IsEnumeration():
			if err := FormatFile(path, enumerationTemplate, e); err != nil {
				log.Fatal(err)
			}
		case e.IsSimple():
			if err := FormatFile(path, simpleTemplate, e); err != nil {
				log.Fatal(err)
			}
		case e.IsDerived():
			if err := FormatFile(path, derivedTemplate, e); err != nil {
				log.Fatal(err)
			}
		default:
			if err := FormatFile(path, complexTemplate, e); err != nil {
				log.Fatal(e.AttrName, err)
			}
		}
	}
}
