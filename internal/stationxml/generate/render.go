package main

import (
	"bytes"
	"go/format"
	"io"
	"os"
)

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
