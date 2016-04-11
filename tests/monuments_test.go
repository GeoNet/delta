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

}
