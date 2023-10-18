package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVolcam(t *testing.T) {

	const res = `[{"id":"raoulisland","mount":"RIMK","view":"01","title":"Raoul Island","latitude":-177.90724,"longitude":-29.267332,"datum":"EPSG:4326","azimuth":270,"height":490,"ground":0,"volcanoes":[{"id":"kermadecislands","title":"Kermadec Islands"}]}]`

	volcs := []Volcam{
		{
			Id:        "raoulisland",
			Mount:     "RIMK",
			View:      "01",
			Title:     "Raoul Island",
			Latitude:  -177.90724,
			Longitude: -29.267332,
			Datum:     "EPSG:4326",
			Azimuth:   270,
			Height:    490,
			Ground:    0,
			Volcanoes: []Volcano{
				{
					Id:    "kermadecislands",
					Title: "Kermadec Islands",
				},
			},
		},
	}

	var buf bytes.Buffer

	if err := Volcams(volcs).Encode(&buf); err != nil {
		t.Fatal(err)
	}

	if s := strings.TrimSpace(buf.String()); s != res {
		t.Errorf("invalid volcam encoding: %s", cmp.Diff(s, res))
	}

}
