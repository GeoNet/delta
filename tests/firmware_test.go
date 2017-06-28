package delta_test

import (
	"sort"
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

func TestFirmware(t *testing.T) {

	var firmwares meta.FirmwareHistoryList
	t.Log("Load firmware history file")
	{
		if err := meta.LoadList("../install/firmware.csv", &firmwares); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for firmware history overlaps")
	{
		installs := make(map[string]meta.FirmwareHistoryList)
		for _, s := range firmwares {
			_, ok := installs[s.Model]
			if ok {
				installs[s.Model] = append(installs[s.Model], s)

			} else {
				installs[s.Model] = meta.FirmwareHistoryList{s}
			}
		}

		var keys []string
		for k, _ := range installs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := installs[k]

			for i, n := 0, len(v); i < n; i++ {
				for j := i + 1; j < n; j++ {
					switch {
					case v[i].Serial != v[j].Serial:
					case v[i].End.Before(v[j].Start):
					case v[i].Start.After(v[j].End):
					case v[i].End.Equal(v[j].Start):
					case v[i].Start.Equal(v[j].End):
					default:
						t.Errorf("firmware %s / %s has overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	var assets meta.AssetList
	t.Log("Load firmware assets file")
	{
		if err := meta.LoadList("../assets/receivers.csv", &assets); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for firmware receiver assets")
	{
		for _, r := range firmwares {
			var found bool
			for _, a := range assets {
				if a.Model != r.Model {
					continue
				}
				if a.Serial != r.Serial {
					continue
				}
				found = true
			}
			if !found {
				t.Errorf("unable to find firmware receiver asset: %s [%s]", r.Model, r.Serial)
			}
		}
	}

	var receivers meta.DeployedReceiverList
	t.Log("Load installed receivers file")
	{
		if err := meta.LoadList("../install/receivers.csv", &receivers); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for latest installed receiver firmware")
	{
		installs := make(map[string]meta.FirmwareHistory)
		for _, s := range firmwares {
			if s.End.Before(time.Now()) {
				continue
			}
			installs[s.Model+"/"+s.Serial] = s
		}
		for _, r := range receivers {
			if r.End.Before(time.Now()) {
				continue
			}
			if _, ok := installs[r.Model+"/"+r.Serial]; !ok {
				t.Errorf("deployed receiver has no current firmware %s / %s at %s between %s and %s", r.Model, r.Serial, r.Mark, r.Start.Format(meta.DateTimeFormat), r.End.Format(meta.DateTimeFormat))
			}
		}
	}
}
