package main

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestToast(t *testing.T) {

	check := Toast{
		XMLName: xml.Name{
			Space: "http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.11",
			Local: "seiscomp",
		},
		NameSpace: "http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.11",
		Version:   "0.11",
		Config: Config{
			ParameterSets: []ParameterSet{
				{
					PublicId: "ParameterSet/trunk/Station/NZ/AUCT/default",
					Created:  "2020-07-27T00:10:55.014699Z",
					ModuleId: "Config/trunk",
					Parameters: []Parameter{
						{
							PublicId: "Parameter/20200727001055.014792.4500",
							Name:     "detecLocid",
							Value:    "40",
						},
						{
							PublicId: "Parameter/20200727001055.014751.4499",
							Name:     "detecStream",
							Value:    "LTZ",
						},
					},
				},
				{
					PublicId: "ParameterSet/trunk/Station/NZ/CHST/default",
					Created:  "2020-07-27T01:55:20.198097Z",
					ModuleId: "Config/trunk",
					Parameters: []Parameter{
						{
							PublicId: "Parameter/20200727015520.198197.621",
							Name:     "detecLocid",
							Value:    "40",
						},
						{
							PublicId: "Parameter/20200727015520.198151.620",
							Name:     "detecStream",
							Value:    "LTZ",
						},
					},
				},
			},
			Module: Module{
				PublicId: "Config/trunk",
				Name:     "trunk",
				Enabled:  true,
				Stations: []Station{
					{
						PublicId:     "Config/trunk/NZ/AUCT",
						NetworkCode:  "NZ",
						StationCode:  "AUCT",
						Enabled:      true,
						AgencyID:     "WEL(GNS_Test)",
						Author:       "trunk",
						CreationTime: "2020-07-27T00:10:55.014503Z",
						Setup: &Setup{
							Name:           "default",
							Enabled:        true,
							ParameterSetID: "ParameterSet/trunk/Station/NZ/AUCT/default",
						},
					},
					{
						PublicId:     "Config/trunk/NZ/CHST",
						NetworkCode:  "NZ",
						StationCode:  "CHST",
						Enabled:      true,
						AgencyID:     "WEL(GNS_Test)",
						Author:       "trunk",
						CreationTime: "2020-07-27T01:55:20.19791Z",
						Setup: &Setup{
							Name:           "default",
							Enabled:        true,
							ParameterSetID: "ParameterSet/trunk/Station/NZ/CHST/default",
						},
					},
				},
			},
		},
	}

	raw, err := os.ReadFile("./testdata/toast.xml")
	if err != nil {
		t.Fatal(err)
	}

	var toast Toast
	if err := toast.Unmarshal(raw); err != nil {
		t.Fatal(err)
	}

	data, err := toast.MarshalIndent("", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(raw, data) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(raw, data))
	}

	if !cmp.Equal(toast, check) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(toast, check))
	}
}
