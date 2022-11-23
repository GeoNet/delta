package meta

import (
	"testing"
	"time"
)

func TestPolarities(t *testing.T) {

	t.Run("check polarities", testListFunc("testdata/polarities.csv", &PolarityList{
		Polarity{
			Station:   "WEL",
			Location:  "10",
			Subsource: "Z",
			Primary:   true,
			Reversed:  true,
			Method:    "compass",
			Span: Span{
				Start: time.Date(2022, 7, 28, 1, 59, 0, 0, time.UTC),
				End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			reversed: "true",
		},
		Polarity{
			Station:  "WEL",
			Location: "11",
			Reversed: true,
			Span: Span{
				Start: time.Date(2022, 7, 28, 1, 59, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			reversed: "true",
		},
		Polarity{
			Station:  "WEL",
			Location: "11",
			Method:   "unknown",
			Span: Span{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			primary:  "false",
			reversed: "false",
		},
	}))
}
