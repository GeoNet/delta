package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var installedRecorderChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for recorder installation overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.InstalledRecorderList)
			for _, s := range set.InstalledRecorders() {
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

						t.Errorf("recorder %s/%s at %-5s has location %-2s overlap between %s and %s",
							v[i].Model, v[i].Serial,
							v[i].Station, v[i].Location,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for invalid installation dates": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, i := range set.InstalledRecorders() {
				if i.End.After(i.Start) {
					continue
				}
				t.Errorf("installed recorder is removed before it has been installed: %s", i.String())
			}
		}
	},

	"check for invalid sensor azimuth": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			for _, i := range set.InstalledRecorders() {
				if i.Orientation.Azimuth < -360.0 || i.Orientation.Azimuth > 360.0 {
					t.Errorf("installed sensor has invalid orientation azimuth: %s [%g]", i.String(), i.Orientation.Azimuth)
				}
			}
		}
	},

	"check for invalid sensor dip": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, i := range set.InstalledRecorders() {
				if i.Orientation.Dip < -90.0 || i.Orientation.Dip > 90.0 {
					t.Errorf("installed sensor has invalid orientation dip: %s [%g]", i.String(), i.Orientation.Dip)
				}
			}
		}
	},

	"check for missing recorder stations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			stations := make(map[string]meta.Station)
			for _, s := range set.Stations() {
				stations[s.Code] = s
			}

			for _, i := range set.InstalledRecorders() {
				if _, ok := stations[i.Station]; !ok {
					t.Errorf("unable to find station: %s", i.Station)
				}
			}

			for _, i := range set.InstalledRecorders() {
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
		}
	},

	"check for missing recorder sites": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			sites := make(map[string]map[string]meta.Site)
			for _, s := range set.Sites() {
				if _, ok := sites[s.Station]; !ok {
					sites[s.Station] = make(map[string]meta.Site)
				}
				sites[s.Station][s.Location] = s
			}

			for _, i := range set.InstalledRecorders() {
				if _, ok := sites[i.Station]; !ok {
					t.Errorf("unable to find sites for station: %s", i.Station)
				}
			}

			for _, i := range set.InstalledRecorders() {
				if s, ok := sites[i.Station]; ok {
					if _, ok := s[i.Location]; !ok {
						t.Errorf("unable to find sites for station/location: %s/%s", i.Station, i.Location)
					}
				}
			}

			for _, i := range set.InstalledRecorders() {
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
		}
	},

	"check for missing assets": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, r := range set.InstalledRecorders() {
				var found bool
				for _, a := range set.Assets() {
					if a.Model != r.DataloggerModel {
						continue
					}
					if a.Serial != r.Serial {
						continue
					}
					found = true
				}
				if found {
					continue
				}
				t.Errorf("unable to find datalogger asset: %s [%s]", r.Model, r.Serial)
			}

		}
	},
}

func TestInstalledRecorders(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range installedRecorderChecks {
		t.Run(k, v(set))
	}
}
