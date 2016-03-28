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

func TestAssets(t *testing.T) {

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

	reference := make(map[string]string)

	for k, v := range files {
		var assets meta.AssetList

		t.Log("Check asset file can be loaded: " + k)
		if err := meta.LoadList(v, &assets); err != nil {
			t.Fatal(err)
		}

		sort.Sort(assets)

		for _, a := range assets {
			if a.AssetNumber != "" {
				if x, ok := reference[a.AssetNumber]; ok {
					t.Error(k + ": Duplicate asset number: " + a.String() + " " + a.AssetNumber + " [" + x + "]")
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

			file, err := ioutil.TempFile(os.TempDir(), "tst")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(file.Name())
			file.Write(buf.Bytes())

			cmd := exec.Command("diff", v, file.Name())
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
	}
}
