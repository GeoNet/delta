package meta

import (
	"testing"
	"time"
)

func TestQuality(t *testing.T) {

	t.Run("check quality", testListFunc("testdata/qualities.csv", &QualityList{
		Quality{
			Station:  "GISS",
			Location: "40",
			Span: Span{
				Start: time.Date(2016, 12, 3, 22, 0, 0, 0, time.UTC),
				End:   time.Date(2016, 12, 14, 19, 0, 0, 0, time.UTC),
			},
			Fault: true,
			fault: "true",
		},
		Quality{
			Station:  "GISS",
			Location: "41",
			Span: Span{
				Start: time.Date(2016, 12, 3, 22, 0, 0, 0, time.UTC),
				End:   time.Date(2016, 12, 14, 19, 0, 0, 0, time.UTC),
			},
			Fault: false,
		},
	}))
}
