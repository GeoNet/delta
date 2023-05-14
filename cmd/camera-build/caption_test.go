package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCameras(t *testing.T) {

	const res = `[{"mount":"CHTB","view":"01","label":"Christchurch East"},{"mount":"CHTB","view":"02","label":"Christchurch North East"}]`

	captions := []Caption{
		{
			Mount: "CHTB",
			View:  "01",
			Label: "Christchurch East",
		},
		{
			Mount: "CHTB",
			View:  "02",
			Label: "Christchurch North East",
		},
	}

	var buf bytes.Buffer

	if err := Captions(captions).Encode(&buf); err != nil {
		t.Fatal(err)
	}

	if s := strings.TrimSpace(buf.String()); s != res {
		t.Errorf("invalid camera encoding: %s", cmp.Diff(s, res))
	}

}
