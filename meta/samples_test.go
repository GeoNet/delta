package meta

import (
	"testing"
	"time"
)

func TestSampleList(t *testing.T) {
	t.Run("check samples", testListFunc("testdata/samples.csv", &SampleList{
		// NA000,MC,Ngauruhoe Volcano,-39.156798,175.631871,2197.0,,WGS84,2003-07-17T00:00:00Z,9999-01-01T00:00:00Z"
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
