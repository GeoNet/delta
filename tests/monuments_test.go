package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var monumentChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for duplicated monuments": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			monuments := set.Monuments()
			for i := 0; i < len(monuments); i++ {
				for j := i + 1; j < len(monuments); j++ {
					if monuments[i].Mark == monuments[j].Mark {
						t.Errorf("monument duplication: %s", monuments[i].Mark)
					}
				}
			}
		}
	},

	"check for monument ground relationships": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			for _, m := range set.Monuments() {
				if m.GroundRelationship > 0.0 {
					t.Errorf("monument has a positive ground relationship: %s [%g]", m.Mark, m.GroundRelationship)
				}
			}
		}
	},

	"check for monument types": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, m := range set.Monuments() {
				switch m.Type {
				case "Shallow Rod / Braced Antenna Mount":
				case "Wyatt/Agnew Drilled-Braced":
				case "Pillar":
				case "Steel Mast":
				case "Unknown":
				default:
					t.Errorf("monument has an unknown type: %s [%s]", m.Mark, m.Type)
				}
			}
		}
	},
}

func TestMonuments(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range monumentChecks {
		t.Run(k, v(set))
	}
}
