package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestGauges(t *testing.T) {

	var gauges meta.GaugeList
	loadListFile(t, "../network/gauges.csv", &gauges)

	t.Run("check for gauge duplication", func(t *testing.T) {
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
	})

	t.Run("check for missing gauge stations ", func(t *testing.T) {
		var list meta.StationList
		loadListFile(t, "../network/stations.csv", &list)

		stas := make(map[string]meta.Station)
		for _, s := range list {
			stas[s.Code] = s
		}
		for _, g := range gauges {
			if _, ok := stas[g.Code]; !ok {
				t.Error("unknown gauge station: " + g.Code)
			}

		}
	})
}
