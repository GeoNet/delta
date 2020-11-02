package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testNetworks = map[string]func([]meta.Network) func(t *testing.T){

	"check for duplicated networks": func(networks []meta.Network) func(t *testing.T) {
		return func(t *testing.T) {

			for i := 0; i < len(networks); i++ {
				for j := i + 1; j < len(networks); j++ {
					if networks[i].Code == networks[j].Code {
						t.Errorf("network duplication: " + networks[i].Code)
					}
				}
			}
		}
	},
}

func TestNetworks(t *testing.T) {
	var networks meta.NetworkList
	loadListFile(t, "../network/networks.csv", &networks)

	for k, fn := range testNetworks {
		t.Run(k, fn(networks))
	}
}
