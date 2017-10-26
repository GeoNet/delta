package delta_test

import (
	"strconv"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestConstituents(t *testing.T) {

	var constituents meta.ConstituentList
	loadListFile(t, "../network/constituents.csv", &constituents)

	t.Run("check for constituent duplications", func(t *testing.T) {
		for i := 0; i < len(constituents); i++ {
			for j := i + 1; j < len(constituents); j++ {
				if constituents[i].Gauge == constituents[j].Gauge && constituents[i].Number == constituents[j].Number {
					t.Error("contituent duplication: " + constituents[i].Gauge + "/" + strconv.Itoa(constituents[i].Number))
				}
			}
		}
	})

	t.Run("check for missing constituent gauges", func(t *testing.T) {
		var list meta.GaugeList
		loadListFile(t, "../network/gauges.csv", &list)

		gauges := make(map[string]meta.Gauge)
		for _, g := range list {
			gauges[g.Code] = g
		}
		for _, c := range constituents {
			if _, ok := gauges[c.Gauge]; !ok {
				t.Error("unknown gauge: " + c.Gauge)
			}
		}
	})
}
