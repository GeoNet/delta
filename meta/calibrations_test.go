package meta

import (
	"testing"
	"time"
)

func TestCalibrationList(t *testing.T) {
	t.Run("check calibrations", testListFunc("testdata/calibrations.csv", &CalibrationList{
		Calibration{
			Install: Install{
				Equipment: Equipment{
					Make:   "Acme",
					Model:  "ACME01",
					Serial: "257",
				},
				Span: Span{
					Start: time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			ScaleFactor:   2000.169 / 2.0,
			ScaleBias:     1.0,
			ScaleAbsolute: 20.0,
			Frequency:     10.0,
			Number:        0,

			number:    "0",
			factor:    "2000.169/2.0",
			bias:      "1.0",
			absolute:  "20",
			frequency: "10.0",
		},
	}))
}
