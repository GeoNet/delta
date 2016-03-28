package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestNetworks(t *testing.T) {

	var networks meta.NetworkList
	t.Log("Load installed sensors file")
	{
		if err := meta.LoadList("../network/networks.csv", &networks); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < len(networks); i++ {
		for j := i + 1; j < len(networks); j++ {
			if networks[i].NetworkCode == networks[j].NetworkCode {
				t.Errorf("network duplication: " + networks[i].NetworkCode)
			}
		}
	}

}
