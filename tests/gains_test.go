package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testGains = map[string]func([]meta.Gain) func(t *testing.T){

	"check for gain installation overlaps": func(installed []meta.Gain) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.GainList)
			for _, s := range installed {
				for _, c := range s.Gains() {
					installs[c.Id()] = append(installs[c.Id()], c)
				}
			}

			for _, v := range installs {
				for i := 0; i < len(v); i++ {
					for j := i + 1; j < len(v); j++ {
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

						t.Errorf("gain %s/%s has component %s overlap between %s and %s",
							v[i].Station, v[i].Location, v[i].Subsource,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},
}

var testGainsSites = map[string]func([]meta.Gain, []meta.Site) func(t *testing.T){
	"check for missing sites": func(installed []meta.Gain, sites []meta.Site) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range installed {
				var found bool
				for _, a := range sites {
					if a.Station != s.Station {
						continue
					}
					if a.Location != s.Location {
						continue
					}
					found = true
				}
				if found {
					continue
				}
				t.Errorf("unable to find gain site: %s [%s]", s.Station, s.Location)
			}

		}
	},
}

func TestGains(t *testing.T) {
	var installed meta.GainList
	loadListFile(t, "../install/gains.csv", &installed)

	for k, fn := range testGains {
		t.Run(k, fn(installed))
	}
}

func TestGains_Sites(t *testing.T) {
	var installed meta.GainList
	loadListFile(t, "../install/gains.csv", &installed)

	var sites meta.SiteList
	loadListFile(t, "../network/sites.csv", &sites)

	for k, fn := range testGainsSites {
		t.Run(k, fn(installed, sites))
	}
}
