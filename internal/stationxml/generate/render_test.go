package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

const testString = "a test string"
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

type Dummy struct{}

func (d Dummy) Render(wr io.Writer, str string) error {
	if _, err := wr.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}

func TestRender(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, testString, Dummy{}); err != nil {
		t.Fatalf("unable to render test string: %v", err)
	}
	if s := buf.String(); s != testString {
		t.Errorf("unable to render test string: found \"%s\", expected \"%s\"", s, testString)
	}
}

func TestRenderFile(t *testing.T) {
	tmp, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	if err := RenderFile(tmp.Name(), testString, Dummy{}); err != nil {
		t.Fatalf("unable to render test string: %v", err)
	}

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("unable to read rendered file: %v", err)
	}

	if s := string(data); s != testString {
		t.Errorf("unable to render test string: found \"%s\", expected \"%s\"", s, testString)
	}
}

func TestFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := Format(&buf, testProgram, Dummy{}); err != nil {
		t.Fatalf("unable to format test program: %v", err)
	}
	if s := buf.String(); s != testFormatted {
		t.Errorf("unable to render test string: found \"%s\", expected \"%s\"", s, testFormatted)
	}
}

func TestFormatFile(t *testing.T) {
	tmp, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	if err := FormatFile(tmp.Name(), testProgram, Dummy{}); err != nil {
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
