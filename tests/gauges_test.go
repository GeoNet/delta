package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestGauges(t *testing.T) {
	var gauges meta.GaugeList

	t.Log("Load installed gauges file")
	if err := meta.LoadList("../network/gauges.csv", &gauges); err != nil {
		t.Fatal(err)
	}

	var stations meta.StationList

	t.Log("Load stations file")
	if err := meta.LoadList("../network/stations.csv", &stations); err != nil {
		t.Fatal(err)
	}

	t.Run("check for guge duplication", func(t *testing.T) {
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

	t.Run("check for missing gauge stations", func(t *testing.T) {

		keys := make(map[string]meta.Station)
		for _, s := range stations {
			keys[s.Code] = s
		}

		for _, g := range gauges {
			if _, ok := keys[g.Code]; !ok {
				t.Error("unknown gauge station: " + g.Code)
			}

		}
	})

}
