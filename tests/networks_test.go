package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestNetworks(t *testing.T) {
	var networks meta.NetworkList

	if err := meta.LoadList("../network/networks.csv", &networks); err != nil {
		t.Fatal(err)
	}

	t.Run("check for duplicate networks", func(t *testing.T) {
		for i := 0; i < len(networks); i++ {
			for j := i + 1; j < len(networks); j++ {
				if networks[i].Code == networks[j].Code {
					t.Errorf("network duplication: " + networks[i].Code)
				}
			}
		}
	})

}
