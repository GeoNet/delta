package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestSensors(t *testing.T) {
	var installed meta.InstalledSensorList
	t.Log("Load installed sensors file")
	{
		if err := meta.LoadList("../install/sensors.csv", &installed); err != nil {
			t.Fatal(err)
		}
		sort.Sort(installed)
	}

	t.Log("Check for missing sensors")
	{
		var sensors meta.AssetList
		if err := meta.LoadList("../assets/sensors.csv", &sensors); err != nil {
			t.Fatal(err)
		}
		sort.Sort(sensors)
		for _, i := range installed {
			n := sort.Search(len(sensors), func(j int) bool { return !sensors[j].Equipment.Less(i.Equipment) })
			if n < 0 || i.Equipment.Less(sensors[n].Equipment) {
				t.Errorf("unable to find sensor: %s", i.String())
			}
		}
	}

	t.Log("Check for sensor installation overlaps")
	{
		installs := make(map[string]meta.InstalledSensorList)
		for _, s := range installed {
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
						t.Errorf("sensor %s/%s at %-5s has location %-2s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].StationCode, v[i].LocationCode, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

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

	for _, i := range installed {
		if i.Orientation.Azimuth < -360.0 || i.Orientation.Azimuth > 360.0 {
			t.Errorf("installed sensor has invalid orientation azimuth: %s [%g]", i.String(), i.Orientation.Azimuth)
		}
		if i.Orientation.Dip < -90.0 || i.Orientation.Dip > 90.0 {
			t.Errorf("installed sensor has invalid orientation dip: %s [%g]", i.String(), i.Orientation.Dip)
		}
	}

}
