package meta

import (
	"testing"
	"time"
)

func TestPreamps(t *testing.T) {

	t.Run("check preamps", testListFunc("testdata/preamps.csv", &PreampList{
		Preamp{
			Station:  "AWAZ",
			Location: "10",
			Span: Span{
				Start: time.Date(2010, 12, 14, 1, 0, 0, 0, time.UTC),
				End:   time.Date(2018, 7, 5, 3, 0, 0, 0, time.UTC),
			},
			Gain: 30,
			gain: "30",
		},
		Preamp{
			Station:  "AWAZ",
			Location: "10",
			Gain:     32,
			Span: Span{
				Start: time.Date(2018, 7, 5, 3, 30, 0, 0, time.UTC),
				End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			gain: "32",
		},
	}))
}
