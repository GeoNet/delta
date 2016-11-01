package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestGauges(t *testing.T) {

	stas := make(map[string]meta.Station)
	{
		var list meta.StationList
		t.Log("Load stations file")
		if err := meta.LoadList("../network/stations.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, s := range list {
			stas[s.Code] = s
		}
	}

	var gauges meta.GaugeList
	t.Log("Load installed gauges file")
	{
		if err := meta.LoadList("../network/gauges.csv", &gauges); err != nil {
			t.Fatal(err)
		}
	}

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

	for _, g := range gauges {
		if _, ok := stas[g.Code]; !ok {
			t.Error("unknown gauge station: " + g.Code)
		}

	}

}
