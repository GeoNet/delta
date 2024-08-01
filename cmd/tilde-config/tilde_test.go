package main

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func mustParse(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &t
}

func TestTilde(t *testing.T) {

	raw, err := os.ReadFile("testdata/tilde.xml")
	if err != nil {
		t.Fatal(err)
	}

	example := Tilde{
		Domains: []Domain{
			{
				Name: "dart",
				Stations: []Station{
					{
						Code: "NZA",
						Latitude: &Float{
							Value: "aaaa",
						},
						Longitude: &Float{
							Value: "bbbb",
						},
						Elevation: &Float{
							Value: "0",
						},
						Sensors: []Sensor{
							{
								Code:  "40",
								Start: mustParse("2019-12-22T00:00:00Z"),
								End:   mustParse("9999-01-01T00:00:00Z"),
								Latitude: &Float{
									Value: "aaaa",
								},
								Longitude: &Float{
									Value: "bbbb",
								},
								Elevation: &Float{
									Value: "0",
								},
								RelativeHeight: &Float{
									Value: "zzzz",
								},
							},
							{
								Code:  "41",
								Start: mustParse("2019-12-22T00:00:00Z"),
								End:   mustParse("9999-01-01T00:00:00Z"),
								Latitude: &Float{
									Value: "aaaa",
								},
								Longitude: &Float{
									Value: "bbbb",
								},
								RelativeHeight: &Float{
									Value: "zzzz",
								},
							},
						},
					},
					{
						Code: "NZB",
						Latitude: &Float{
							Value: "cccc",
						},
						Longitude: &Float{
							Value: "dddd",
						},
						Elevation: &Float{
							Value: "0",
						},
						Sensors: []Sensor{
							{
								Code:  "40",
								Start: mustParse("2019-12-22T00:00:00Z"),
								End:   mustParse("9999-01-01T00:00:00Z"),
								Latitude: &Float{
									Value: "cccc",
								},
								Longitude: &Float{
									Value: "dddd",
								},
								Elevation: &Float{
									Value: "0",
								},
								RelativeHeight: &Float{
									Value: "yyyy",
								},
							},
							{
								Code:  "41",
								Start: mustParse("2019-12-22T00:00:00Z"),
								End:   mustParse("9999-01-01T00:00:00Z"),
								Latitude: &Float{
									Value: "cccc",
								},
								Longitude: &Float{
									Value: "dddd",
								},
								RelativeHeight: &Float{
									Value: "yyyy",
								},
							},
						},
					},
				},
			},
			{
				Name: "fits",
			},
		},
	}

	var buf bytes.Buffer

	if err := example.MarshalIndent(&buf, "", "  "); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf.Bytes(), raw) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(buf.String(), string(raw)))
	}
}
