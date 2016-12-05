package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestMonuments(t *testing.T) {
	var monuments meta.MonumentList

	if err := meta.LoadList("../network/monuments.csv", &monuments); err != nil {
		t.Fatal(err)
	}

	t.Run("check for duplicate monuments", func(t *testing.T) {
		for i := 0; i < len(monuments); i++ {
			for j := i + 1; j < len(monuments); j++ {
				if monuments[i].Mark == monuments[j].Mark {
					t.Errorf("monument duplication: " + monuments[i].Mark)
				}
			}
		}
	})

	t.Run("check monument details", func(t *testing.T) {
		for _, m := range monuments {
			if m.GroundRelationship > 0.0 {
				t.Errorf("positive monuments ground relationship: %s [%g]", m.Mark, m.GroundRelationship)
			}

			switch m.Type {
			case "Shallow Rod / Braced Antenna Mount":
			case "Wyatt/Agnew Drilled-Braced":
			case "Pillar":
			case "Steel Mast":
			case "Unknown":
			default:
				t.Errorf("unknown monument type: %s [%s]", m.Mark, m.Type)
			}
		}
	})
}
