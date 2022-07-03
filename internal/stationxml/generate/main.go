package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"aqwari.net/xml/xsdgen"
)

func download(wr io.Writer, url string, insecure bool) error {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, //nolint:gosec // needed until fdsn.org passed checks
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = io.Copy(wr, resp.Body); err != nil {
		return err
	}

	return nil
}

func main() {

	var name string
	flag.StringVar(&name, "name", "stationxml", "package name")

	var remote string
	flag.StringVar(&remote, "schema", "https://www.fdsn.org/xml/station/fdsn-station-1.2.xsd", "schema service endpoint to download from")

	var base string
	flag.StringVar(&base, "base", ".", "base directory")

	var version string
	flag.StringVar(&version, "version", "v1.2", "optional schema version to add to directory path")

	var output string
	flag.StringVar(&output, "output", "stationxml.go", "output file")

	var insecure bool
	flag.BoolVar(&insecure, "insecure", false, "whether the remote site has certificate issues, use with caution")

	flag.Parse()

	file, err := os.CreateTemp(".", "schema")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	if err := download(file, remote, insecure); err != nil {
		log.Fatal(err)
	}

	// where to store the output file
	path := filepath.Join(base, version, output)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Fatal(err)
	}

	var cfg xsdgen.Config
	cfg.Option(
		xsdgen.PackageName(name),
		xsdgen.UseFieldNames(),
	)

	code, err := cfg.GenSource(file.Name())
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(path, code, 0644); err != nil {
		log.Fatal(err)
	}
}
