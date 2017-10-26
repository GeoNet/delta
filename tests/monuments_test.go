package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestMonuments(t *testing.T) {

	var monuments meta.MonumentList
	loadListFile(t, "../network/monuments.csv", &monuments)

	t.Run("check for monument duplication", func(t *testing.T) {
		for i := 0; i < len(monuments); i++ {
			for j := i + 1; j < len(monuments); j++ {
				if monuments[i].Mark == monuments[j].Mark {
					t.Errorf("monument duplication: " + monuments[i].Mark)
				}
			}
		}
	})

	t.Run("check for monument ground relationships", func(t *testing.T) {
		for _, m := range monuments {
			if m.GroundRelationship > 0.0 {
				t.Errorf("monument has a positive ground relationship: %s [%g]", m.Mark, m.GroundRelationship)
			}
		}
	})

	t.Run("check for monument types", func(t *testing.T) {
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
	})
}
