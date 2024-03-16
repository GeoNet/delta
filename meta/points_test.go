package meta

import (
	"testing"
	"time"
)

func TestPointlList(t *testing.T) {
	t.Run("check channels", testListFunc("testdata/points.csv", &PointList{
		Point{
			Sample:   "NA000",
			Location: "MC01",
			Position: Position{
				Latitude:  -39.156798,
				Longitude: 175.631871,
				Elevation: 2197.0,
				Depth:     0.0,
				Datum:     "WGS84",

				latitude:  "-39.156798",
				longitude: "175.631871",
				elevation: "2197.0",
				depth:     "",
			},
			Span: Span{
				Start: time.Date(2003, 7, 17, 0, 0, 0, 0, time.UTC),
				End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			Survey: "External GPS Device",
		},
		Point{
			Sample:   "NA001",
			Location: "MC01",
			Position: Position{
				Latitude:  -39.1562653309,
				Longitude: 175.631464065,
				Elevation: 2280.0,
				Depth:     0.0,
				Datum:     "WGS84",

				latitude:  "-39.1562653309",
				longitude: "175.631464065",
				elevation: "2280",
				depth:     "",
			},
			Span: Span{
				Start: time.Date(1978, 7, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			Survey: "External GPS Device",
		},
		Point{
			Sample:   "NA002",
			Location: "MC01",
			Position: Position{
				Latitude:  -39.156415,
				Longitude: 175.634722,
				Elevation: 2280.0,
				Depth:     0.0,
				Datum:     "WGS84",

				latitude:  "-39.156415",
				longitude: "175.634722",
				elevation: "2280",
				depth:     "",
			},
			Span: Span{
				Start: time.Date(1973, 11, 27, 0, 0, 0, 0, time.UTC),
				End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			Survey: "External GPS Device",
		},
	}))
}
