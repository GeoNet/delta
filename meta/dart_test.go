package meta

import (
	"testing"
	"time"
)

func TestDarts(t *testing.T) {

	t.Run("check preamps", testListFunc("testdata/darts.csv", &DartList{
		Dart{
			Station:       "NZA",
			Pid:           "5501002",
			WmoIdentifier: "NZ41",

			Span: Span{
				Start: time.Date(2019, 12, 21, 17, 0, 0, 0, time.UTC),
				End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		Dart{
			Station:       "NZB",
			Pid:           "5501003",
			WmoIdentifier: "NZ42",
			Span: Span{
				Start: time.Date(2020, 9, 18, 0, 0, 0, 0, time.UTC),
				End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}))
}
