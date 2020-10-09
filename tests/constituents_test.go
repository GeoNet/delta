package delta_test

import (
	"strconv"
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testConstituents = map[string]func([]meta.Constituent) func(t *testing.T){
	"check for gauge duplication": func(constituents []meta.Constituent) func(t *testing.T) {
		return func(t *testing.T) {
			for i := 0; i < len(constituents); i++ {
				for j := i + 1; j < len(constituents); j++ {
					if constituents[i].Gauge == constituents[j].Gauge && constituents[i].Number == constituents[j].Number {
						t.Error("contituent duplication: " + constituents[i].Gauge + "/" + strconv.Itoa(constituents[i].Number))
					}
				}
			}
		}
	},
}

var testConstituentsGauges = map[string]func([]meta.Constituent, []meta.Gauge) func(t *testing.T){
	"check for missing constituent gauges": func(constituents []meta.Constituent, list []meta.Gauge) func(t *testing.T) {
		return func(t *testing.T) {
			gauges := make(map[string]meta.Gauge)
			for _, g := range list {
				gauges[g.Code] = g
			}
			for _, c := range constituents {
				if _, ok := gauges[c.Gauge]; !ok {
					t.Error("unknown gauge: " + c.Gauge)
				}
			}
		}
	},
}

func TestConstituents(t *testing.T) {
	var constituents meta.ConstituentList
	loadListFile(t, "../environment/constituents.csv", &constituents)

	for k, fn := range testConstituents {
		t.Run(k, fn(constituents))
	}
}

func TestConstituents_Gauges(t *testing.T) {
	var constituents meta.ConstituentList
	loadListFile(t, "../environment/constituents.csv", &constituents)

	var gauges meta.GaugeList
	loadListFile(t, "../environment/gauges.csv", &gauges)

	for k, fn := range testConstituentsGauges {
		t.Run(k, fn(constituents, gauges))
	}

}
