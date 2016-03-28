package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestStations(t *testing.T) {

	var stations meta.StationList
	t.Log("Load installed sensors file")
	{
		if err := meta.LoadList("../network/stations.csv", &stations); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < len(stations); i++ {
		for j := i + 1; j < len(stations); j++ {
			if stations[i].Code == stations[j].Code {
				t.Errorf("station duplication: " + stations[i].Code)
			}
		}
	}

}
