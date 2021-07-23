package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testMeasurements = map[string]func([]meta.Measurement) func(t *testing.T){
	"check for duplicated measurements": func(measurements []meta.Measurement) func(t *testing.T) {
		return func(t *testing.T) {

			for i := 0; i < len(measurements); i++ {
				for j := i + 1; j < len(measurements); j++ {
					if measurements[i].Location != measurements[j].Location {
						continue
					}
					if measurements[i].Name != measurements[j].Name {
						continue
					}
					if measurements[i].Sensor != measurements[j].Sensor {
						continue
					}
					if measurements[i].Type != measurements[j].Type {
						continue
					}
					t.Errorf("measurement duplication: " + measurements[i].Name + "/" + measurements[j].Name)
				}
			}
		}
	},
}

func TestMeasurements(t *testing.T) {

	var measurements meta.MeasurementList
	loadListFile(t, "../fits/measurements.csv", &measurements)

	for k, fn := range testMeasurements {
		t.Run(k, fn(measurements))
	}
}
