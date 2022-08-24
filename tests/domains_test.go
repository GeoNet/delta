package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var testDomains = map[string]func(*meta.Set) func(t *testing.T){

	"check for valid domains": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			stations := make(map[string]meta.Station)
			for _, s := range set.Stations() {
				stations[s.Code] = s
			}

			networks := make(map[string]meta.Network)
			for _, n := range set.Networks() {
				networks[n.Code] = n
			}

			for _, d := range set.Domains() {
				if _, ok := networks[d.Network]; ok {
					continue
				}
				t.Errorf("unable to find network: %s", d.Network)
			}

			for _, d := range set.Domains() {
				if _, ok := networks[d.Domain]; ok {
					continue
				}
				t.Errorf("unable to find domain: %s", d.Domain)

			}

			for _, d := range set.Domains() {
				if _, ok := stations[d.Station]; ok {
					continue
				}
				t.Errorf("unable to find station: %s", d.Station)
			}
		}
	},
}

func TestDomains(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, fn := range testDomains {
		t.Run(k, fn(set))
	}
}
