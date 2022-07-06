package main

import (
	"bytes"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// FileName returns a snake-case version of the struct based on it's name, it will attempt to
// replace camel case runs with underscore breaks as well as full uppercase names.
func FileName(name string) string {
	var label string
	var multi bool
	for _, r := range filepath.Base(name) {
		switch l := unicode.ToLower(r); {
		case l == unicode.ToUpper(l):
			multi = false
		case unicode.IsUpper(r) && len(label) > 0 && !multi:
			label += "_" + string(l)
			multi = true
		case unicode.IsUpper(r):
			label += string(l)
			multi = true
		default:
			label += string(l)
			multi = false
		}
	}

	// handle uppercase name bleeding into Type
	switch {
	case label == "type":
	case strings.HasSuffix(label, "_type"):
	case !strings.HasSuffix(label, "type"):
	default:
		label = strings.TrimSuffix(label, "type") + "_type"
	}

	// building a golang file
	return filepath.Join(filepath.Dir(name), label+".go")
}

// Renderer is an interface to describe writing a rendered go template into a writer.
type Renderer interface {
	Render(io.Writer, string) error
}

// Render writes a rendered go template into a writer.
func Render(wr io.Writer, tmpl string, renderer Renderer) error {

	var bytes bytes.Buffer
	if err := renderer.Render(&bytes, tmpl); err != nil {
		return err
	}

	if _, err := wr.Write(bytes.Bytes()); err != nil {
		return err
	}

	return nil
}

// RenderFile writes a rendered template into a file.
func RenderFile(name string, tmpl string, renderer Renderer) error {

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	return Render(file, tmpl, renderer)
}

// Format writes a rendered go template into a writer.
func Format(wr io.Writer, tmpl string, renderer Renderer) error {

	var bytes bytes.Buffer
	if err := renderer.Render(&bytes, tmpl); err != nil {
		return err
	}

	formatted, err := format.Source(bytes.Bytes())
	if err != nil {
		return err
	}

	if _, err := wr.Write(formatted); err != nil {
		return err
	}

	return nil
}

// FormatFile writes a rendered go template into a file.
func FormatFile(name string, tmpl string, renderer Renderer) error {

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	return Format(file, tmpl, renderer)
}
