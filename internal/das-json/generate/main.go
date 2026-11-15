package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func FetchSchema(url string, file string) error {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false}, //nolint:gosec // needed until fdsn.org passed checks
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close() // ignore body close errors
	}()
	out, err := os.Create(file)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil

}

func main() {

	var name string
	flag.StringVar(&name, "name", "dasjson", "package name")

	var base string
	flag.StringVar(&base, "base", "dasjson.go", "file name for base go code")

	var version float64
	flag.Float64Var(&version, "version", 0.0, "schema version encoded into module")

	var datetime string
	flag.StringVar(&datetime, "datetime", "date_time.go", "file name for extra date time go code")

	var channels string
	flag.StringVar(&channels, "channels", "Channels.go", "file name for replacement channels go code")

	var future bool
	flag.BoolVar(&future, "future", false, "output dates in the future")

	var format string
	flag.StringVar(&format, "format", "2006-01-02T15:04:05Z", "provide date time format to encode as")

	var input string
	flag.StringVar(&input, "input", "", "input schema file")

	var remote string
	flag.StringVar(&remote, "schema", "", "schema service endpoint to download from")

	var insecure bool
	flag.BoolVar(&insecure, "insecure", false, "whether the remote site has certificate issues, use with caution")

	var ns string
	flag.StringVar(&ns, "ns", "https://raw.githubusercontent.com/FDSN/DAS-metadata/refs/heads/main/schema/DAS-Metadata.v2.0.schema.json", "schema namespace to process")

	var output string
	flag.StringVar(&output, "output", "output", "output dir")

	var null string
	flag.StringVar(&null, "null", "Created,FrequencyStart,FrequencyEnd,FrequencyDBVariation", "related to groups some variables should be pointers")

	var schemaFile string
	flag.StringVar(&schemaFile, "file", "/home/tbenson/git/delta/das-metadata-test/internal/das-json/generate/schema.json", "Schema File Save location")

	flag.Parse()

	// err := FetchSchema(ns, schemaFile)

	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failure fetching schema file", err.Error())
	// 	os.Exit(1)
	// }

	schema, err := ReadInputFile(schemaFile, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failure reading Schema file", err.Error())
		os.Exit(1)
	}

	g := New(schema)
	// fmt.Print(*g, "\n\n")
	// fmt.Print(*g.schemas[0], "\n\n")

	err = g.CreateTypes()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failure generating structs: ", err)
		os.Exit(1)
	}

	if output != "" {
		// w, err = os.Create(*&output)
		err := os.MkdirAll(output, os.ModePerm)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating output directory: ", err)
			return
		}
	}
	// a := g.schemas[0].JSONKey
	log.Printf("rendering %s => %s", g.schemas[0].ID06, output)
	Output(output, g, name)

	settings := Datetime{
		Package: name,
		Format:  format,
		Future:  future,
	}

	path := filepath.Join(output, datetime)
	log.Printf("rendering datetime => %s", path)
	if err := settings.RenderFile(path); err != nil {
		log.Fatalf("unable to parse date time type: %v", err)
	}
	// path = filepath.Join(output, channels)
	// chan_settings := Channels{
	// 	Package: name,
	// }
	// log.Printf("rendering channels => %s", path)
	// if err := chan_settings.RenderFile(path); err != nil {
	// 	log.Fatalf("unable to parse: %v", err)
	// }

}
