package delta_test

import (
	"strings"
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var classChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for duplicated classes": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			classes := set.Classes()

			for i := 0; i < len(classes); i++ {
				for j := i + 1; j < len(classes); j++ {
					if classes[i].Station == classes[j].Station {
						t.Errorf("class site duplication: %s <=> %s", classes[i].Station, classes[j].Station)
					}
				}
			}
		}
	},

	"check for missing class stations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			for _, c := range set.Classes() {
				if strings.HasPrefix(c.Station, "#") {
					continue
				}
				if _, ok := set.Station(c.Station); !ok {
					t.Errorf("class station missing: %s", c.Station)
				}
			}
		}
	},

	"check for invalid class settings": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, c := range set.Classes() {
				switch c.SiteClass {
				case "A", "B", "C", "D", "E":
				default:
					t.Errorf("class invalid site class for station %s: %s", c.Station, c.SiteClass)
				}
				switch c.Vs30Quality {
				case "Q3", "Q2", "Q1":
				default:
					t.Errorf("class invalid Vs30 quality for station %s: %s", c.Station, c.Vs30Quality)
				}
				switch c.TsiteQuality {
				case "I", "Q3", "Q2", "Q1":
				default:
					t.Errorf("class invalid Tsite quality for station %s: %s", c.Station, c.TsiteQuality)
				}
				switch c.TsiteMethod {
				case "I", "Ms", "Mw", "Mn", "Mu", "Ma":
				default:
					t.Errorf("class invalid Tsite method for station %s: %s", c.Station, c.TsiteMethod)
				}

				switch c.DepthQuality {
				case "Q3", "Q2", "Q1":
				default:
					t.Errorf("class invalid depth quality for station %s: %s", c.Station, c.DepthQuality)
				}
			}
		}
	},

	"check for missing class citations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, c := range set.Classes() {
				for _, r := range c.Citations {
					if _, ok := set.Citation(r); !ok {
						t.Errorf("class unknown citation for station %s: %q", c.Station, r)
					}
				}
			}
		}
	},
}

func TestClasses(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range classChecks {
		t.Run(k, v(set))
	}
}
