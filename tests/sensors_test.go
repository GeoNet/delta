package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestSensors(t *testing.T) {
	var installed meta.InstalledSensorList
	loadListFile(t, "../install/sensors.csv", &installed)

	t.Run("check for missing sensors", func(t *testing.T) {
		var sensors meta.AssetList
		loadListFile(t, "../assets/sensors.csv", &sensors)
		sort.Sort(sensors)
		for _, i := range installed {
			n := sort.Search(len(sensors), func(j int) bool { return !sensors[j].Equipment.Less(i.Equipment) })
			if n < 0 || i.Equipment.Less(sensors[n].Equipment) {
				t.Errorf("unable to find sensor: %s", i.String())
			}
		}
	})

	t.Run("check for sensor installation overlaps", func(t *testing.T) {
		installs := make(map[string]meta.InstalledSensorList)
		for _, s := range installed {
			if _, ok := installs[s.Model]; !ok {
				installs[s.Model] = meta.InstalledSensorList{}
			}
			installs[s.Model] = append(installs[s.Model], s)
		}

		for _, v := range installs {
			for i := 0; i < len(v); i++ {
				for j := i + 1; j < len(v); j++ {
					if v[i].Serial != v[j].Serial {
						continue
					}
					if v[i].End.Before(v[j].Start) {
						continue
					}
					if v[i].Start.After(v[j].End) {
						continue
					}
					if v[i].End.Equal(v[j].Start) {
						continue
					}
					if v[i].Start.Equal(v[j].End) {
						continue
					}

					t.Errorf("sensor %s/%s at %-5s has location %-2s overlap between %s and %s",
						v[i].Model, v[i].Serial,
						v[i].Station, v[i].Location,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("check for missing sensor stations", func(t *testing.T) {
		var list meta.StationList
		loadListFile(t, "../network/stations.csv", &list)

		stations := make(map[string]meta.Station)
		for _, s := range list {
			stations[s.Code] = s
		}

		for _, i := range installed {
			if _, ok := stations[i.Station]; !ok {
				t.Errorf("unable to find station: %s", i.Station)
			}
		}

		for _, i := range installed {
			if s, ok := stations[i.Station]; ok {
				if i.Start.Before(s.Start) {
					t.Logf("warning: installed sensor before station has been opened: %s: %s (%s %s)",
						i.String(), i.Start.String(), s.Code, s.Start.String())
				}
				if i.End.After(s.End) {
					t.Logf("warning: installed sensor after station has been closed: %s: %s (%s %s)",
						i.String(), i.End.String(), s.Code, s.End.String())
				}
			}
		}
	})

	t.Run("check for missing sensor sites", func(t *testing.T) {
		var list meta.SiteList
		loadListFile(t, "../network/sites.csv", &list)

		sites := make(map[string]map[string]meta.Site)
		for _, s := range list {
			if _, ok := sites[s.Station]; !ok {
				sites[s.Station] = make(map[string]meta.Site)
			}
			sites[s.Station][s.Location] = s
		}

		for _, i := range installed {
			if _, ok := sites[i.Station]; !ok {
				t.Errorf("unable to find sites for station: %s", i.Station)
			}
		}

		for _, i := range installed {
			if s, ok := sites[i.Station]; ok {
				if _, ok := s[i.Location]; !ok {
					t.Errorf("unable to find sites for station/location: %s/%s", i.Station, i.Location)
				}
			}
		}

		for _, i := range installed {
			if s, ok := sites[i.Station]; ok {
				if l, ok := s[i.Location]; ok {
					if i.Start.Before(l.Start) {
						t.Logf("warning: installed sensor before site has been opened: %s: %s (%s/%s %s)",
							i.String(), i.Start.String(), l.Station, l.Location, l.Start.String())
					}
					if i.End.After(l.End) {
						t.Logf("warning: installed sensor after site has been closed: %s: %s (%s/%s %s)",
							i.String(), i.End.String(), l.Station, l.Location, l.End.String())
					}
				}
			}
		}
	})

	t.Run("check for invalid sensor azimuth", func(t *testing.T) {
		for _, i := range installed {
			if i.Orientation.Azimuth < -360.0 || i.Orientation.Azimuth > 360.0 {
				t.Errorf("installed sensor has invalid orientation azimuth: %s [%g]", i.String(), i.Orientation.Azimuth)
			}
		}
	})

	t.Run("check for invalid sensor dip", func(t *testing.T) {
		for _, i := range installed {
			if i.Orientation.Dip < -90.0 || i.Orientation.Dip > 90.0 {
				t.Errorf("installed sensor has invalid orientation dip: %s [%g]", i.String(), i.Orientation.Dip)
			}
		}
	})

	t.Run("check for missing sensor assets", func(t *testing.T) {
		var assets meta.AssetList
		loadListFile(t, "../assets/sensors.csv", &assets)
		for _, s := range installed {
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
			if found {
				continue
			}
			t.Errorf("unable to find sensor asset: %s [%s]", s.Model, s.Serial)
		}
	})

}
