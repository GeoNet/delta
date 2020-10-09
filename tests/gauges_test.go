package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testGauges = map[string]func([]meta.Gauge) func(t *testing.T){
	"check for gauge duplication": func(gauges []meta.Gauge) func(t *testing.T) {
		return func(t *testing.T) {
			for i := 0; i < len(gauges); i++ {
				for j := i + 1; j < len(gauges); j++ {
					if gauges[i].Code == gauges[j].Code {
						t.Errorf("gauge code duplication: " + gauges[i].Code)
					}
					if gauges[i].Number == gauges[j].Number {
						t.Errorf("gauge number duplication: " + gauges[i].Code)
					}
				}
			}
		}
	},
}

var testGaugesStations = map[string]func([]meta.Gauge, []meta.Station) func(t *testing.T){
	"check for missing gauge stations": func(gauges []meta.Gauge, stations []meta.Station) func(t *testing.T) {
		return func(t *testing.T) {

			stas := make(map[string]meta.Station)
			for _, s := range stations {
				stas[s.Code] = s
			}

			for _, g := range gauges {
				if _, ok := stas[g.Code]; !ok {
					t.Error("unknown gauge station: " + g.Code)
				}

			}
		}
	},
}

func TestGauges(t *testing.T) {

	var gauges meta.GaugeList
	loadListFile(t, "../environment/gauges.csv", &gauges)

	for k, fn := range testGauges {
		t.Run(k, fn(gauges))
	}
}

func TestGauges_Stations(t *testing.T) {

	var gauges meta.GaugeList
	loadListFile(t, "../environment/gauges.csv", &gauges)

	var stations meta.StationList
	loadListFile(t, "../network/stations.csv", &stations)

	for k, fn := range testGaugesStations {
		t.Run(k, fn(gauges, stations))
	}

}
