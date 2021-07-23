package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testLocations = map[string]func([]meta.Location) func(t *testing.T){
	"check for duplicated sites": func(sites []meta.Location) func(t *testing.T) {
		return func(t *testing.T) {

			for i := 0; i < len(sites); i++ {
				for j := i + 1; j < len(sites); j++ {
					if sites[i].Code != sites[j].Code {
						continue
					}
					t.Errorf("site duplication: " + sites[i].Code)
				}
			}
		}
	},
}

func TestLocations(t *testing.T) {

	var locations meta.LocationList
	loadListFile(t, "../fits/locations.csv", &locations)

	for k, fn := range testLocations {
		t.Run(k, fn(locations))
	}
}
