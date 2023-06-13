package delta_test

import (
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var siteChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicated sites": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			sites := set.Sites()
			for i := 0; i < len(sites); i++ {
				for j := i + 1; j < len(sites); j++ {
					if sites[i].Station == sites[j].Station && sites[i].Location == sites[j].Location {
						t.Errorf("site duplication: " + sites[i].Station + "/" + sites[i].Location)
					}
				}
			}
		}
	},

	"check for duplicated station sites": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			stas := make(map[string]meta.Station)
			for _, s := range set.Stations() {
				stas[s.Code] = s
			}
			for _, c := range set.Sites() {
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

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range siteChecks {
		t.Run(k, v(set))
	}
}
