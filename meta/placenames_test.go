package meta

import (
	"path/filepath"
	"testing"
)

func TestPlacenames(t *testing.T) {

	places, err := LoadPlacenames(filepath.Join("testdata", filepath.Base(PlacenamesFile)))
	if err != nil {
		t.Fatal(err)
	}

	checks := []struct {
		lat  float64
		lon  float64
		desc string
	}{
		{-36.85, 174.767, "within 5 km of Auckland"},
		{-36.733, 174.683, "15 km north-west of Auckland"},
		{-41.283, 174.767, "382 km north-east of Ashburton"},
		{-43.533, 172.633, "82 km north-east of Ashburton"},
		{-45.883, 170.5, "242 km south-west of Ashburton"},
	}

	for _, c := range checks {
		if d := PlacenameList(places).Description(c.lat, c.lon); d != c.desc {
			t.Errorf("invalid placename description for %g/%g, expected %q, got %q", c.lat, c.lon, c.desc, d)
		}
	}
}

func TestPlacenameList(t *testing.T) {

	t.Run("check placenames", testListFunc("testdata/placenames.csv", &PlacenameList{
		Placename{
			Name:      "Ashburton",
			Latitude:  -43.9,
			Longitude: 171.75,
			Level:     1,

			latitude:  "-43.9",
			longitude: "171.75",
		},
		Placename{
			Name:      "Auckland",
			Latitude:  -36.85,
			Longitude: 174.767,
			Level:     0,

			latitude:  "-36.85",
			longitude: "174.767",
		},
	}))
}
