package delta_test

import (
	"strconv"
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var constituentChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for gauge duplication": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			constituents := set.Constituents()
			for i := 0; i < len(constituents); i++ {
				for j := i + 1; j < len(constituents); j++ {
					if constituents[i].Gauge != constituents[j].Gauge {
						continue
					}
					if constituents[i].Number != constituents[j].Number {
						continue
					}
					if !constituents[i].Start.Equal(constituents[j].Start) {
						continue
					}

					t.Error("constituent duplication: " + constituents[i].Gauge + "/" + strconv.Itoa(constituents[i].Number))
				}
			}
		}
	},
	"check for missing constituent gauges": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			gauges := make(map[string]meta.Gauge)
			for _, g := range set.Gauges() {
				gauges[g.Code] = g
			}
			for _, c := range set.Constituents() {
				if _, ok := gauges[c.Gauge]; !ok {
					t.Error("unknown gauge: " + c.Gauge)
				}
			}
		}
	},
}

func TestConstituents(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range constituentChecks {
		t.Run(k, v(set))
	}
}
