package delta_test

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestAssets(t *testing.T) {

	files := map[string]string{
		"antennas":    "../equipment/antennas.csv",
		"cameras":     "../equipment/cameras.csv",
		"dataloggers": "../equipment/dataloggers.csv",
		"metsensors":  "../equipment/metsensors.csv",
		"radomes":     "../equipment/radomes.csv",
		"receivers":   "../equipment/receivers.csv",
		"recorders":   "../equipment/recorders.csv",
		"sensors":     "../equipment/sensors.csv",
	}

	reference := make(map[string]string)

	for k, v := range files {
		var assets meta.AssetList

		t.Log("Check asset file can be loaded: " + k)
		if err := meta.LoadList(v, &assets); err != nil {
			t.Fatal(err)
		}

		for _, a := range assets {
			if a.AssetNumber != "" {
				if _, ok := reference[a.AssetNumber]; ok {
					t.Error(k + ": Duplicate asset number: " + a.String() + " " + a.AssetNumber)
				}
				reference[a.AssetNumber] = a.String()
			}
		}

		t.Log("Check asset file consistency: " + k)
		raw, err := ioutil.ReadFile(v)
		if err != nil {
			t.Fatal(err)
		}

		var buf bytes.Buffer
		if err := csv.NewWriter(&buf).WriteAll(meta.EncodeList(assets)); err != nil {
			t.Fatal(err)
		}

		if string(raw) != buf.String() {
			t.Error(k + ": Assets file mismatch: " + v)
		}
	}
}
