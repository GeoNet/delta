package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestSensors(t *testing.T) {
	var sensors meta.InstalledSensorList

	if err := meta.LoadList("../install/sensors.csv", &sensors); err != nil {
		t.Fatal(err)
	}

	var assets meta.AssetList

	if err := meta.LoadList("../assets/sensors.csv", &assets); err != nil {
		t.Fatal(err)
	}

	t.Run("Check for missing sensors", func(t *testing.T) {
		for _, i := range sensors {
			n := sort.Search(len(assets), func(j int) bool { return !assets[j].Equipment.Less(i.Equipment) })
			if n < 0 || i.Equipment.Less(assets[n].Equipment) {
				t.Errorf("unable to find sensor: %s", i.String())
			}
		}
	})

	t.Log("Check for sensor installation overlaps", func(t *testing.T) {
		installs := make(map[string]meta.InstalledSensorList)
		for _, s := range sensors {
			_, ok := installs[s.Model]
			if ok {
				installs[s.Model] = append(installs[s.Model], s)

			} else {
				installs[s.Model] = meta.InstalledSensorList{s}
			}
		}

		var keys []string
		for k, _ := range installs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			for i, v, n := 0, installs[k], len(installs[k]); i < n; i++ {
				for j := i + 1; j < n; j++ {
					if v[i].Serial != v[j].Serial {
						continue
					}
					if v[i].End.Before(v[j].Start) || v[i].Start.After(v[j].End) {
						continue
					}
					if v[i].End.Equal(v[j].Start) || v[i].Start.Equal(v[j].End) {
						continue
					}
					t.Errorf("sensor %s/%s at %-5s has location %-2s overlap between %s and %s",
						v[i].Model, v[i].Serial, v[i].Station, v[i].Location,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	/*
		t.Log("Check for missing sensor stations")
		{
			var stations meta.StationList

			if err := meta.LoadList("../network/stations.csv", &stations); err != nil {
				t.Fatal(err)
			}

			sort.Sort(stations)

			for _, i := range installed {
				n := sort.Search(len(stations), func(j int) bool { return !(stations[j].Code < i.StationCode) })
				if n < 0 || i.StationCode != stations[n].Code {
					t.Errorf("unable to find station: %s", i.StationCode)
				} else {
					if i.Start.Before(stations[n].Start) {
						t.Errorf("installed sensor before station has been opened: %s: %s (%s %s)", i.String(), i.Start.String(), stations[n].Code, stations[n].Start.String())
					}
					if i.End.After(stations[n].End) {
						t.Errorf("installed sensor after station has been closed: %s: %s (%s %s)", i.String(), i.End.String(), stations[n].Code, stations[n].End.String())
					}
				}
			}
		}

		t.Log("Check for missing sensor sites")
		{
			var sites meta.SiteList

			if err := meta.LoadList("../network/sites.csv", &sites); err != nil {
				t.Fatal(err)
			}

			sort.Sort(sites)

			for _, i := range installed {
				n := sort.Search(len(sites), func(j int) bool {
					if sites[j].StationCode > i.StationCode {
						return true
					}
					if sites[j].StationCode < i.StationCode {
						return false
					}
					return !(sites[j].LocationCode < i.LocationCode)
				})
				if n < 0 || i.StationCode != sites[n].StationCode || i.LocationCode != sites[n].LocationCode {
					t.Errorf("unable to find site: %s/%s", i.StationCode, i.LocationCode)
				} else {
					if i.Start.Before(sites[n].Start) {
						t.Errorf("installed sensor before site has been opened: %s: %s (%s/%s %s)", i.String(), i.Start.String(), sites[n].StationCode, sites[n].LocationCode, sites[n].Start.String())
					}
					if i.End.After(sites[n].End) {
						t.Errorf("installed sensor after site has been closed: %s: %s (%s/%s %s)", i.String(), i.End.String(), sites[n].StationCode, sites[n].LocationCode, sites[n].End.String())
					}
				}
			}
		}
	*/

	t.Run("Check for sensor installations", func(t *testing.T) {
		for _, i := range sensors {
			if i.Orientation.Azimuth < -360.0 || i.Orientation.Azimuth > 360.0 {
				t.Errorf("installed sensor has invalid orientation azimuth: %s [%g]", i.String(), i.Orientation.Azimuth)
			}
			if i.Orientation.Dip < -90.0 || i.Orientation.Dip > 90.0 {
				t.Errorf("installed sensor has invalid orientation dip: %s [%g]", i.String(), i.Orientation.Dip)
			}
		}
	})

	t.Run("Check for sensor assets", func(t *testing.T) {
		for _, s := range sensors {
			var found bool
			for _, a := range assets {
				if a.Model != s.Model {
					continue
				}
				if a.Serial != s.Serial {
					continue
				}
				found = true
			}
			if !found {
				t.Errorf("unable to find sensor asset: %s [%s]", s.Model, s.Serial)
			}
		}
	})

}
