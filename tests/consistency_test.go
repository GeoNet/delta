package delta_test

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/GeoNet/delta/meta"
)

func loadListFile(t *testing.T, path string, list meta.List) {
	if err := meta.LoadList(path, list); err != nil {
		t.Fatalf("unable to load list file %s: %v", path, err)
	}
}

var testConsistency = map[string]func(path string, list meta.List) func(t *testing.T){

	"check file consistency": func(path string, list meta.List) func(t *testing.T) {
		return func(t *testing.T) {

			raw, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("unable to read %s file: %v", path, err)
			}

			var buf bytes.Buffer
			if err := csv.NewWriter(&buf).WriteAll(meta.EncodeList(list)); err != nil {
				t.Fatalf("unable to decode %s file: %v", path, err)
			}

			if string(raw) != buf.String() {
				t.Errorf("unexpected %s content -got/+exp\n%s", filepath.Base(path), cmp.Diff(string(raw), buf.String()))
			}
		}
	},
}

func TestConsistency(t *testing.T) {

	files := map[string]struct {
		f string
		l meta.List
	}{
		"connections":  {f: "../install/connections.csv", l: &meta.ConnectionList{}},
		"cameras":      {f: "../install/cameras.csv", l: &meta.InstalledCameraList{}},
		"dataloggers":  {f: "../install/dataloggers.csv", l: &meta.DeployedDataloggerList{}},
		"metsensors":   {f: "../install/metsensors.csv", l: &meta.InstalledMetSensorList{}},
		"radomes":      {f: "../install/radomes.csv", l: &meta.InstalledRadomeList{}},
		"receivers":    {f: "../install/receivers.csv", l: &meta.DeployedReceiverList{}},
		"recorders":    {f: "../install/recorders.csv", l: &meta.InstalledRecorderList{}},
		"sensors":      {f: "../install/sensors.csv", l: &meta.InstalledSensorList{}},
		"firmware":     {f: "../install/firmware.csv", l: &meta.FirmwareHistoryList{}},
		"streams":      {f: "../install/streams.csv", l: &meta.StreamList{}},
		"networks":     {f: "../network/networks.csv", l: &meta.NetworkList{}},
		"stations":     {f: "../network/stations.csv", l: &meta.StationList{}},
		"sites":        {f: "../network/sites.csv", l: &meta.SiteList{}},
		"marks":        {f: "../network/marks.csv", l: &meta.MarkList{}},
		"mounts":       {f: "../network/mounts.csv", l: &meta.MountList{}},
		"views":        {f: "../network/views.csv", l: &meta.ViewList{}},
		"gauges":       {f: "../environment/gauges.csv", l: &meta.GaugeList{}},
		"constituents": {f: "../environment/constituents.csv", l: &meta.ConstituentList{}},
	}

	for f, v := range files {

		loadListFile(t, v.f, v.l)

		for k, fn := range testConsistency {
			t.Run(k+": "+f, fn(v.f, v.l))
		}
	}
}
