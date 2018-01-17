package delta_test

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func loadListFile(t *testing.T, path string, list meta.List) {
	if err := meta.LoadList(path, list); err != nil {
		t.Fatalf("unable to load list file %s: %v", path, err)
	}
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
		"gauges":       {f: "../environment/gauges.csv", l: &meta.GaugeList{}},
		"constituents": {f: "../environment/constituents.csv", l: &meta.ConstituentList{}},
	}

	for k, v := range files {
		t.Run("check file consistency: "+k, func(t *testing.T) {
			loadListFile(t, v.f, v.l)

			sort.Sort(v.l)

			raw, err := ioutil.ReadFile(v.f)
			if err != nil {
				t.Fatalf("unable to read %s asset file: %v", k, err)
			}

			var buf bytes.Buffer
			if err := csv.NewWriter(&buf).WriteAll(meta.EncodeList(v.l)); err != nil {
				t.Fatalf("unable to decode %s asset file: %v", k, err)
			}

			if string(raw) != buf.String() {
				t.Errorf("%s: **** csv file mismatch **** : ", v.f)

				file, err := ioutil.TempFile(os.TempDir(), "tst")
				if err != nil {
					t.Fatalf("unable to create temporary file: %v", err)
				}
				defer os.Remove(file.Name())
				file.Write(buf.Bytes())

				cmd := exec.Command("diff", "-c", v.f, file.Name())
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					t.Fatalf("unable to create diff: %v", err)
				}
				err = cmd.Start()
				if err != nil {
					t.Fatalf("unable to start diff: %v", err)
				}
				defer cmd.Wait()
				diff, err := ioutil.ReadAll(stdout)
				if err != nil {
					t.Fatalf("unable to read diff: %v", err)
				}

				t.Error(string(diff))
			}
		})
	}
}
