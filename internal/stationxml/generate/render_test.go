package main

import (
	"bytes"
	"os"
	"testing"
)

const testProgram = `package main
const test = "test"
func test() string {
	return test
}
`

const testFormatted = `package main

const test = "test"

func test() string {
	return test
}
`

type mirror struct{}

func (d mirror) Render() ([]byte, error) {
	return []byte(testProgram), nil
}

func TestRender(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, mirror{}); err != nil {
		t.Fatalf("unable to render test string: %v", err)
	}
	if s := buf.String(); s != testFormatted {
		t.Errorf("unable to render test program: found \"%s\", expected \"%s\"", s, testFormatted)
	}
}

func TestRenderFile(t *testing.T) {
	tmp, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	if err := RenderFile(tmp.Name(), mirror{}); err != nil {
		t.Fatalf("unable to render test program: %v", err)
	}

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("unable to read rendered file: %v", err)
	}

	if s := string(data); s != testFormatted {
		t.Errorf("unable to render test string: found \"%s\", expected \"%s\"", s, testFormatted)
	}
}
