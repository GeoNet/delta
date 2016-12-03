package delta_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestAssets(t *testing.T) {
	var assets meta.AssetList

	files := map[string]string{
		"antennas":    "../assets/antennas.csv",
		"cameras":     "../assets/cameras.csv",
		"dataloggers": "../assets/dataloggers.csv",
		"metsensors":  "../assets/metsensors.csv",
		"radomes":     "../assets/radomes.csv",
		"receivers":   "../assets/receivers.csv",
		"recorders":   "../assets/recorders.csv",
		"sensors":     "../assets/sensors.csv",
	}

	for k, v := range files {
		var a meta.AssetList
		t.Logf("Load %s assets file", k)
		if err := meta.LoadList(v, &a); err != nil {
			t.Fatal(err)
		}
		assets = append(assets, a...)
	}

	sort.Sort(assets)

	t.Run("check duplicate assets", func(t *testing.T) {
		reference := make(map[string]string)

		for _, a := range assets {
			if a.Number == "" {
				continue
			}
			if x, ok := reference[a.Number]; ok {
				t.Error("duplicate asset number: " + a.String() + " " + a.Number + " [" + x + "]")
			}
			reference[a.Number] = a.String()
		}

		for i := 0; i < len(assets); i++ {
			for j := i + 1; j < len(assets); j++ {
				if assets[i].Model != assets[j].Model {
					continue
				}
				if assets[i].Serial != assets[j].Serial {
					continue
				}
				t.Errorf("equipment duplication: " + strings.Join([]string{assets[i].Model, assets[i].Serial}, " "))
			}
		}
	})
}
