package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var testStations = map[string]func(set *meta.Set) func(t *testing.T){
	"check for duplicated stations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			stations := set.Stations()

			for i := 0; i < len(stations); i++ {
				for j := i + 1; j < len(stations); j++ {
					if stations[i].Code != stations[j].Code {
						continue
					}
					if stations[i].Network != stations[j].Network {
						continue
					}
					t.Errorf("station duplication: %s/%s", stations[i].Code, stations[i].Network)
				}
			}
		}
	},
	"check for missing networks": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			nets := make(map[string]meta.Network)
			for _, n := range set.Networks() {
				nets[n.Code] = n
			}
			for _, s := range set.Stations() {
				if _, ok := nets[s.Network]; ok {
					continue
				}
				t.Logf("warning: missing network %s: %s", s.Code, s.Network)
			}
		}
	},
}

func TestStations(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, fn := range testStations {
		t.Run(k, fn(set))
	}
}
