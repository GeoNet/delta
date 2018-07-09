package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestNetworks(t *testing.T) {

	var networks meta.NetworkList
	loadListFile(t, "../network/networks.csv", &networks)

	t.Run("check for network duplication", func(t *testing.T) {
		for i := 0; i < len(networks); i++ {
			for j := i + 1; j < len(networks); j++ {
				if networks[i].Code == networks[j].Code {
					t.Errorf("network duplication: " + networks[i].Code)
				}
			}
		}
	})
}
