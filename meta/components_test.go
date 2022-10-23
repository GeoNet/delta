package meta

import (
	"testing"
)

func TestComponentList(t *testing.T) {
	t.Run("check components", testListFunc("testdata/components.csv", &ComponentList{
		Component{
			Make:         "Guralp",
			Model:        "Fortis",
			Type:         "Accelerometer",
			Number:       0,
			Subsource:    "Z",
			Dip:          -90.0,
			Azimuth:      0.0,
			Types:        "G",
			Response:     "sensor_guralp_fortis_response",
			SamplingRate: 10,

			number:       "0",
			dip:          "-90",
			azimuth:      "0",
			samplingRate: "10",
		},
		Component{
			Make:      "Guralp",
			Model:     "Fortis",
			Type:      "Accelerometer",
			Number:    1,
			Subsource: "N",
			Dip:       0.0,
			Azimuth:   0.0,
			Types:     "G",
			Response:  "sensor_guralp_fortis_response",

			number:  "1",
			dip:     "0",
			azimuth: "0",
		},
		Component{
			Make:         "Guralp",
			Model:        "Fortis",
			Type:         "Accelerometer",
			Number:       2,
			Subsource:    "E",
			Dip:          0.0,
			Azimuth:      90.0,
			Types:        "G",
			SamplingRate: 1 / 10,
			Response:     "sensor_guralp_fortis_response",

			number:       "2",
			dip:          "0",
			azimuth:      "90",
			samplingRate: "-10",
		},
	}))
}
