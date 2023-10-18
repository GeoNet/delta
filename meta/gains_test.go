package meta

import (
	"testing"
	"time"
)

func TestGainList(t *testing.T) {
	t.Run("check gains", testListFunc("testdata/gains.csv", &GainList{
		Gain{
			Span: Span{
				Start: time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			Scale: Scale{
				Factor: 1298.169,
				Bias:   11865.556,

				factor: "1298.169",
				bias:   "11865.556",
			},
			Station:     "SBAM",
			Location:    "50",
			Sublocation: "01",
			Subsource:   "XZ",
		},
	}))
}
