package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var placenameChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for placename duplication": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			placenames := set.Placenames()
			for i := 0; i < len(placenames); i++ {
				for j := i + 1; j < len(placenames); j++ {
					if placenames[i].Name != placenames[j].Name {
						continue
					}
					t.Errorf("placename duplication: %s", placenames[i].Name)
				}
			}
		}
	},
	"check for placename latitude longitudes": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, p := range set.Placenames() {
				if p.Latitude < -90.0 || p.Latitude > 90.0 {
					t.Errorf("placename latitude problem: %s (%g)", p.Name, p.Latitude)
				}
				if p.Longitude < -180.0 || p.Longitude > 180.0 {
					t.Errorf("placename longitude problem: %s (%g)", p.Name, p.Longitude)
				}
			}
		}
	},
	"check for placename levels": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, p := range set.Placenames() {
				if p.Level < 0 || p.Level > 3 {
					t.Errorf("placename level problem: %s (%d)", p.Name, p.Level)
				}
			}
		}
	},
}

func TestPlacenames(t *testing.T) {
	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range placenameChecks {
		t.Run(k, v(set))
	}
}
