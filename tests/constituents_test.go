package delta_test

import (
	"strconv"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestConstituents(t *testing.T) {
	var constituents meta.ConstituentList

	t.Log("Load installed constituents file")
	if err := meta.LoadList("../network/constituents.csv", &constituents); err != nil {
		t.Fatal(err)
	}

	var gauges meta.GaugeList

	t.Log("Load installed gauges file")
	if err := meta.LoadList("../network/gauges.csv", &gauges); err != nil {
		t.Fatal(err)
	}

	t.Run("check for constituent duplication", func(t *testing.T) {
		for i := 0; i < len(constituents); i++ {
			for j := i + 1; j < len(constituents); j++ {
				if constituents[i].Gauge != constituents[j].Gauge {
					continue
				}
				if constituents[i].Number != constituents[j].Number {
					continue
				}
				t.Error("contituent duplication: " + constituents[i].Gauge + "/" + strconv.Itoa(constituents[i].Number))
			}
		}
	})

	t.Run("check for missing constituent gauges", func(t *testing.T) {
		list := make(map[string]meta.Gauge)
		for _, g := range gauges {
			list[g.Code] = g
		}
		for _, c := range constituents {
			if _, ok := list[c.Gauge]; !ok {
				t.Error("unknown gauge: " + c.Gauge)
			}

		}
	})
}
