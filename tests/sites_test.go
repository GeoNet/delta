package delta_test

import (
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

var testSites = map[string]func([]meta.Site) func(t *testing.T){
	"check for duplicated sites": func(sites []meta.Site) func(t *testing.T) {
		return func(t *testing.T) {

			for i := 0; i < len(sites); i++ {
				for j := i + 1; j < len(sites); j++ {
					if sites[i].Station == sites[j].Station && sites[i].Location == sites[j].Location {
						t.Errorf("site duplication: " + sites[i].Station + "/" + sites[i].Location)
					}
				}
			}
		}
	},
}

var testSites_Stations = map[string]func([]meta.Site, []meta.Station) func(t *testing.T){
	"check for duplicated sites": func(sites []meta.Site, stations []meta.Station) func(t *testing.T) {
		return func(t *testing.T) {
			stas := make(map[string]meta.Station)
			for _, s := range stations {
				stas[s.Code] = s
			}
			for _, c := range sites {
				if s, ok := stas[c.Station]; ok {
					switch {
					case c.Start.Before(s.Start):
						t.Log("warning: site start mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.Start.String(),
							"before",
							s.Start.String(),
						}, " "))
					case s.End.Before(time.Now()) && c.End.After(s.End):
						t.Log("warning: site end mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.End.String(),
							"after",
							s.End.String(),
						}, " "))
					}
				}
			}
		}
	},
}

func TestSites(t *testing.T) {

	var sites meta.SiteList
	loadListFile(t, "../network/sites.csv", &sites)

	for k, fn := range testSites {
		t.Run(k, fn(sites))
	}
}

func TestSites_Stations(t *testing.T) {

	var sites meta.SiteList
	loadListFile(t, "../network/sites.csv", &sites)

	var stations meta.StationList
	loadListFile(t, "../network/stations.csv", &stations)

	for k, fn := range testSites_Stations {
		t.Run(k, fn(sites, stations))
	}
}
