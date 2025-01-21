package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var gaugeChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for gauge duplication": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			gauges := set.Gauges()
			for i := 0; i < len(gauges); i++ {
				for j := i + 1; j < len(gauges); j++ {
					if gauges[i].Code != gauges[j].Code {
						continue
					}
					if !gauges[i].Start.Equal(gauges[j].Start) {
						continue
					}
					t.Errorf("gauge code duplication: %s", gauges[i].Code)
				}
			}
			for i := 0; i < len(gauges); i++ {
				for j := i + 1; j < len(gauges); j++ {
					if gauges[i].Number == "" {
						continue
					}
					if gauges[i].Number != gauges[j].Number {
						continue
					}
					if !gauges[i].Start.Equal(gauges[j].Start) {
						continue
					}
					t.Errorf("gauge number duplication: %s", gauges[i].Code)
				}
			}
		}
	},
	"check for missing gauge stations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			stas := make(map[string]meta.Station)
			for _, s := range set.Stations() {
				stas[s.Code] = s
			}

			for _, g := range set.Gauges() {
				if _, ok := stas[g.Code]; !ok {
					t.Errorf("unknown gauge station: %s", g.Code)
				}

			}
		}
	},
}

func TestGauges(t *testing.T) {
	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range gaugeChecks {
		t.Run(k, v(set))
	}
}
