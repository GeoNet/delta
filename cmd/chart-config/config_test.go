package main

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfig(t *testing.T) {

	config := Config([]Stream{
		{
			Srcname:            "NZ_CAW_10_EHZ",
			NetworkCode:        "NZ",
			StationCode:        "CAW",
			LocationCode:       "10",
			ChannelCode:        "EHZ",
			StationName:        "Cannon Point",
			InternalNetwork:    "WL",
			NetworkDescription: "Wellington regional seismic network",
			Latitude:           -41.107194232,
			Longitude:          175.066438523,
			SamplingPeriod:     10000000,
			Sensitivity:        167772160,
			Gain:               1,
			InputUnits:         "m/s",
			OutputUnits:        "count",
		},
	})

	raw, err := os.ReadFile("./testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}

	data, err := config.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(raw, data) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(raw, data))
	}
}
