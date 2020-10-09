package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testMonuments = map[string]func([]meta.Monument) func(t *testing.T){

	"check for duplicated monuments": func(monuments []meta.Monument) func(t *testing.T) {
		return func(t *testing.T) {

			for i := 0; i < len(monuments); i++ {
				for j := i + 1; j < len(monuments); j++ {
					if monuments[i].Mark == monuments[j].Mark {
						t.Errorf("monument duplication: " + monuments[i].Mark)
					}
				}
			}
		}
	},

	"check for monument ground relationships": func(monuments []meta.Monument) func(t *testing.T) {
		return func(t *testing.T) {

			for _, m := range monuments {
				if m.GroundRelationship > 0.0 {
					t.Errorf("monument has a positive ground relationship: %s [%g]", m.Mark, m.GroundRelationship)
				}
			}
		}
	},

	"check for monument types": func(monuments []meta.Monument) func(t *testing.T) {
		return func(t *testing.T) {
			for _, m := range monuments {
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
	var monuments meta.MonumentList
	loadListFile(t, "../network/monuments.csv", &monuments)

	for k, fn := range testMonuments {
		t.Run(k, fn(monuments))
	}
}
