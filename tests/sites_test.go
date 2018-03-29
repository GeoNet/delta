package delta_test

import (
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

func TestSites(t *testing.T) {

	var sites meta.SiteList
	loadListFile(t, "../network/sites.csv", &sites)

	stas := make(map[string]meta.Station)
	t.Run("load stations file", func(t *testing.T) {
		var list meta.StationList
		loadListFile(t, "../network/stations.csv", &list)
		for _, s := range list {
			stas[s.Code] = s
		}
	})

	t.Run("check for duplicated sites", func(t *testing.T) {
		for i := 0; i < len(sites); i++ {
			for j := i + 1; j < len(sites); j++ {
				if sites[i].Station == sites[j].Station && sites[i].Location == sites[j].Location {
					t.Errorf("site duplication: " + sites[i].Station + "/" + sites[i].Location)
				}
			}
		}
	})

	t.Run("check for invalid site dates", func(t *testing.T) {
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
	})

}
