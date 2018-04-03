package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestStations(t *testing.T) {

	var stations meta.StationList
	loadListFile(t, "../network/stations.csv", &stations)

	t.Run("check for duplicated stations", func(t *testing.T) {
		for i := 0; i < len(stations); i++ {
			for j := i + 1; j < len(stations); j++ {
				if stations[i].Code == stations[j].Code {
					t.Errorf("station duplication: " + stations[i].Code)
				}
			}
		}
	})

}
