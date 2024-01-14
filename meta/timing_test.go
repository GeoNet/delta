package meta

import (
	"testing"
	"time"
)

func TestTimings(t *testing.T) {

	t.Run("check timings", testListFunc("testdata/timings.csv", &TimingList{
		Timing{
			Station:    "NZB",
			Location:   "42",
			Correction: -24 * time.Hour,
			Span: Span{
				Start: time.Date(2023, 12, 7, 14, 15, 0, 0, time.UTC),
				End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			correction: "-24h",
		},
	}))
}
