package main

import (
	"go/format"
	"io"
	"os"
)

// Renderer is an interface which as single function, Render which implements
// the ability for a type to convert itself into go code. This code is passed
// through the standard go formatter and available for storage as required.
type Renderer interface {
	Render() ([]byte, error)
}

// RenderFile applies a Renderer and stores its go code output into the given file path.
func RenderFile(path string, renderer Renderer) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return Render(file, renderer)
}

// Render applies a Renderer and passes its go code output into the Writer provided.
// Render takes a Renderer and stores its go code into the given file path.
func Render(wr io.Writer, renderer Renderer) error {
	data, err := renderer.Render()
	if err != nil {
		return err
	}

	formatted, err := format.Source(data)
	if err != nil {
		return err
	}

	if _, err := wr.Write(formatted); err != nil {
		return err
	}

	return nil
}
