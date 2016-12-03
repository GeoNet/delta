package delta_test

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

// TestConsistency will ensure that after reading each csv file will be the same on write. The system application
// diff will be called on files that show differences.
func TestConsistency(t *testing.T) {

	files := map[string]struct {
		p string
		l meta.List
	}{
		"connections":  {"../install/connections.csv", &meta.ConnectionList{}},
		"cameras":      {"../install/cameras.csv", &meta.InstalledCameraList{}},
		"dataloggers":  {"../install/dataloggers.csv", &meta.DeployedDataloggerList{}},
		"metsensors":   {"../install/metsensors.csv", &meta.InstalledMetSensorList{}},
		"radomes":      {"../install/radomes.csv", &meta.InstalledRadomeList{}},
		"receivers":    {"../install/receivers.csv", &meta.DeployedReceiverList{}},
		"recorders":    {"../install/recorders.csv", &meta.InstalledRecorderList{}},
		"sensors":      {"../install/sensors.csv", &meta.InstalledSensorList{}},
		"firmware":     {"../install/firmware.csv", &meta.FirmwareHistoryList{}},
		"streams":      {"../install/streams.csv", &meta.StreamList{}},
		"networks":     {"../network/networks.csv", &meta.NetworkList{}},
		"stations":     {"../network/stations.csv", &meta.StationList{}},
		"sites":        {"../network/sites.csv", &meta.SiteList{}},
		"marks":        {"../network/marks.csv", &meta.MarkList{}},
		"mounts":       {"../network/mounts.csv", &meta.MountList{}},
		"gauges":       {"../network/gauges.csv", &meta.GaugeList{}},
		"constituents": {"../network/constituents.csv", &meta.ConstituentList{}},

		"antenna assets":    {"../assets/antennas.csv", &meta.AssetList{}},
		"camera assets":     {"../assets/cameras.csv", &meta.AssetList{}},
		"datalogger assets": {"../assets/dataloggers.csv", &meta.AssetList{}},
		"metsensor assets":  {"../assets/metsensors.csv", &meta.AssetList{}},
		"radome assets":     {"../assets/radomes.csv", &meta.AssetList{}},
		"receiver assets":   {"../assets/receivers.csv", &meta.AssetList{}},
		"recorder assets":   {"../assets/recorders.csv", &meta.AssetList{}},
		"sensor assets":     {"../assets/sensors.csv", &meta.AssetList{}},
	}

	for k, v := range files {
		t.Run(fmt.Sprintf("%s: %s", k, filepath.Base(v.p)), func(t *testing.T) {

			raw, err := ioutil.ReadFile(v.p)
			if err != nil {
				t.Fatal(err)
			}

			if err := meta.LoadList(v.p, v.l); err != nil {
				t.Fatal(err)
			}

			sort.Sort(v.l)

			var buf bytes.Buffer
			if err := csv.NewWriter(&buf).WriteAll(meta.EncodeList(v.l)); err != nil {
				t.Fatal(err)
			}

			if string(raw) != buf.String() {
				t.Error(k + ": **** csv file mismatch **** : " + v.p)

				file, err := ioutil.TempFile(os.TempDir(), "tst")
				if err != nil {
					t.Fatal(err)
				}
				defer os.Remove(file.Name())
				file.Write(buf.Bytes())

				cmd := exec.Command("diff", "-c", v.p, file.Name())
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					t.Fatal(err)
				}
				err = cmd.Start()
				if err != nil {
					t.Fatal(err)
				}
				defer cmd.Wait()
				diff, err := ioutil.ReadAll(stdout)
				if err != nil {
					t.Fatal(err)
				}
				t.Error(string(diff))
			}
		})

	}
}
