package meta

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func testHeaderFunc(path string, list List) func(t *testing.T) {
	return func(t *testing.T) {
		// this is the output that is expected
		res, err := MarshalList(list)
		if err != nil {
			t.Fatal(err)
		}

		// load the updated file back into list
		if err := LoadList(path, list); err != nil {
			t.Fatal(err)
		}

		// rebuild the list output, they should be the same
		check, err := MarshalList(list)
		if err != nil {
			t.Fatal(err)
		}

		if string(res) != string(check) {
			t.Errorf("unexpected %s content -got/+exp\n%s", path, cmp.Diff(res, check))
		}
	}
}

func TestHeaderList(t *testing.T) {

	t.Run("check reordered headers", testHeaderFunc("testdata/reorder.csv", &SampleList{
		Sample{
			Reference: Reference{
				Code:    "NA000",
				Network: "MC",
				Name:    "Ngauruhoe Volcano",
			},
			Position: Position{
				Latitude:  -39.156798,
				Longitude: 175.631871,
				Elevation: 2197.0,
				Depth:     0,
				Datum:     "WGS84",

				latitude:  "-39.156798",
				longitude: "175.631871",
				elevation: "2197.0",
				depth:     "",
			},
			Span: Span{
				Start: time.Date(2003, time.July, 17, 0, 0, 0, 0, time.UTC),
				End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		Sample{
			Reference: Reference{
				Code:    "NA001",
				Network: "MC",
				Name:    "Ngauruhoe crater inner rim north-west fumarole field",
			},
			Position: Position{
				Latitude:  -39.1562653309,
				Longitude: 175.631464065,
				Elevation: 2280.0,
				Depth:     0,
				Datum:     "WGS84",

				latitude:  "-39.1562653309",
				longitude: "175.631464065",
				elevation: "2280",
				depth:     "",
			},
			Span: Span{
				Start: time.Date(1978, time.July, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}))

	t.Run("check extra headers", testHeaderFunc("testdata/extra.csv", &SampleList{
		Sample{
			Reference: Reference{
				Code:    "NA000",
				Network: "MC",
				Name:    "Ngauruhoe Volcano",
			},
			Position: Position{
				Latitude:  -39.156798,
				Longitude: 175.631871,
				Elevation: 2197.0,
				Depth:     0,
				Datum:     "WGS84",

				latitude:  "-39.156798",
				longitude: "175.631871",
				elevation: "2197.0",
				depth:     "",
			},
			Span: Span{
				Start: time.Date(2003, time.July, 17, 0, 0, 0, 0, time.UTC),
				End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		Sample{
			Reference: Reference{
				Code:    "NA001",
				Network: "MC",
				Name:    "Ngauruhoe crater inner rim north-west fumarole field",
			},
			Position: Position{
				Latitude:  -39.1562653309,
				Longitude: 175.631464065,
				Elevation: 2280.0,
				Depth:     0,
				Datum:     "WGS84",

				latitude:  "-39.1562653309",
				longitude: "175.631464065",
				elevation: "2280",
				depth:     "",
			},
			Span: Span{
				Start: time.Date(1978, time.July, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}))
}
