package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testStations = map[string]func([]meta.Station) func(t *testing.T){
	"check for duplicated stations": func(stations []meta.Station) func(t *testing.T) {
		return func(t *testing.T) {
			for i := 0; i < len(stations); i++ {
				for j := i + 1; j < len(stations); j++ {
					if stations[i].Code == stations[j].Code {
						t.Errorf("station duplication: " + stations[i].Code)
					}
				}
			}
		}
	},
}

var testStationsNetworks = map[string]func([]meta.Station, []meta.Network) func(t *testing.T){
	"check for missing networks": func(stations []meta.Station, networks []meta.Network) func(t *testing.T) {
		return func(t *testing.T) {
			nets := make(map[string]meta.Network)
			for _, n := range networks {
				nets[n.Code] = n
			}
			for _, s := range stations {
				if _, ok := nets[s.Network]; !ok {
					t.Logf("warning: missing network %s: %s", s.Code, s.Network)
				}
			}
		}
	},
}

func TestStations(t *testing.T) {

	var stations meta.StationList
	loadListFile(t, "../network/stations.csv", &stations)

	for k, fn := range testStations {
		t.Run(k, fn(stations))
	}
}

func TestStations_Networks(t *testing.T) {

	var stations meta.StationList
	loadListFile(t, "../network/stations.csv", &stations)

	var networks meta.NetworkList
	loadListFile(t, "../network/networks.csv", &networks)

	for k, fn := range testStationsNetworks {
		t.Run(k, fn(stations, networks))
	}
}
