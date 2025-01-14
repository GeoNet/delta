package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var networkChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for duplicated networks": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			networks := set.Networks()
			for i := 0; i < len(networks); i++ {
				for j := i + 1; j < len(networks); j++ {
					if networks[i].Code == networks[j].Code {
						t.Errorf("network duplication: %s", networks[i].Code)
					}
				}
			}
		}
	},
}

func TestNetworks(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range networkChecks {
		t.Run(k, v(set))
	}
}
