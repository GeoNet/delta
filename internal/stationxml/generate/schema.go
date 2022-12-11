package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io"
	"net/http"
	"os"

	"aqwari.net/xml/xsd"
)

// TODO: go1.18 has this as a stdlib builtin
func cut(s, sep []byte) (before, after []byte, found bool) {
	if i := bytes.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, nil, false
}

// escape xml like constructs inside documentation
func recomment(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	var tagFront = []byte("<xs:documentation>")
	var tagBack = []byte("</xs:documentation>")

	for {
		front, back, ok := cut(data, tagFront)
		if _, err := buf.Write(front); err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		if _, err := buf.Write(tagFront); err != nil {
			return nil, err
		}
		inside, next, ok := cut(back, tagBack)

		replaced := bytes.ReplaceAll(bytes.ReplaceAll(inside, []byte("<"), []byte("&lt;")), []byte(">"), []byte("&gt;"))
		if _, err := buf.Write(replaced); err != nil {
			return nil, err
		}

		// needs to always end in any case
		if _, err := buf.Write(tagBack); err != nil {
			return nil, err
		}

		if !ok {
			break
		}

		data = next
	}

	return buf.Bytes(), nil
}

// Schemas is a wrapper for a slice of xsd.Schema values
type Schemas struct {
	Namespace string
	Schema    []xsd.Schema
	Pointers  []string

	docs map[xml.Name]string
}

func NewSchemas(namespace string, pointers ...string) *Schemas {
	return &Schemas{
		Namespace: namespace,
		Pointers:  pointers,
		docs:      make(map[xml.Name]string),
	}
}

// Read decodes an FDSN StationXML schema from a byte slice.
func (s *Schemas) Read(data []byte) error {

	updated, err := recomment(data)
	if err != nil {
		return err
	}

	schemas, err := xsd.Parse(updated)
	if err != nil {
		return err
	}

	s.Schema = append(s.Schema, schemas...)

	return nil
}

// Download and decode a remote FDSN StationXML schemas
func (s *Schemas) Download(url string, insecure bool) error {
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

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, resp.Body); err != nil {
		return err
	}

	return s.Read(buf.Bytes())
}

// ReadFile decodes an FDSN StationXML schema from a file.
func (s *Schemas) ReadFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return s.Read(raw)
}
