package delta_test

import (
	"strconv"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestConstituents(t *testing.T) {

	gauges := make(map[string]meta.Gauge)
	{
		t.Log("Load installed gauges file")
		var list meta.GaugeList
		if err := meta.LoadList("../network/gauges.csv", &list); err != nil {
			t.Fatal(err)
		}

		for _, g := range list {
			gauges[g.Code] = g
		}
	}

	var constituents meta.ConstituentList
	if err := meta.LoadList("../network/constituents.csv", &constituents); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(constituents); i++ {
		for j := i + 1; j < len(constituents); j++ {
			if constituents[i].Gauge == constituents[j].Gauge && constituents[i].Number == constituents[j].Number {
				t.Error("contituent duplication: " + constituents[i].Gauge + "/" + strconv.Itoa(constituents[i].Number))
			}
		}
	}

	for _, c := range constituents {
		if _, ok := gauges[c.Gauge]; !ok {
			t.Error("unknown gauge: " + c.Gauge)
		}

	}
}
