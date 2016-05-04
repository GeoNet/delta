package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestMonuments(t *testing.T) {

	var monuments meta.MonumentList
	t.Log("Load network monuments file")
	{
		if err := meta.LoadList("../network/monuments.csv", &monuments); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < len(monuments); i++ {
		for j := i + 1; j < len(monuments); j++ {
			if monuments[i].MarkCode == monuments[j].MarkCode {
				t.Errorf("monument duplication: " + monuments[i].MarkCode)
			}
		}
	}

	for _, m := range monuments {
		if m.GroundRelationship > 0.0 {
			t.Errorf("positive monuments ground relationship: %s [%g]", m.MarkCode, m.GroundRelationship)
		}

		switch m.MonumentType {
		case "Shallow Rod / Braced Antenna Mount":
		case "Wyatt/Agnew Drilled-Braced":
		case "Pillar":
		case "Steel Mast":
		case "Unknown":
		default:
			t.Errorf("unknown monument type: %s [%s]", m.MarkCode, m.MonumentType)
		}
	}
}
