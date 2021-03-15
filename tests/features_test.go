package delta_test

import (
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

var testFeatures = map[string]func([]meta.Feature) func(t *testing.T){
	"check for duplicated site features": func(features []meta.Feature) func(t *testing.T) {
		return func(t *testing.T) {

			for i := 0; i < len(features); i++ {
				for j := i + 1; j < len(features); j++ {
					if features[i].Station != features[j].Station {
						continue
					}
					if features[i].Location == features[j].Location {
						continue
					}
					if features[i].End.Before(features[j].Start) {
						continue
					}
					t.Errorf("site feature overlap: " + features[i].Station + "/" + features[i].Location)
				}
			}
		}
	},
}

var testFeatures_Stations = map[string]func([]meta.Feature, []meta.Station) func(t *testing.T){
	"check for duplicated features": func(features []meta.Feature, stations []meta.Station) func(t *testing.T) {
		return func(t *testing.T) {
			stas := make(map[string]meta.Station)
			for _, s := range stations {
				stas[s.Code] = s
			}
			for _, c := range features {
				if _, ok := stas[c.Station]; !ok {
					t.Error("error: unable to find feature station: " + c.Station)
				}
			}
			for _, c := range features {
				if s, ok := stas[c.Station]; ok {
					switch {
					case c.Start.Before(s.Start):
						t.Log("warning: feature start mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.Start.String(),
							"before",
							s.Start.String(),
						}, " "))
					case s.End.Before(time.Now()) && c.End.After(s.End):
						t.Log("warning: feature end mismatch: " + strings.Join([]string{
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

var testFeatures_Sites = map[string]func([]meta.Feature, []meta.Site) func(t *testing.T){
	"check for duplicated feature sites": func(features []meta.Feature, sites []meta.Site) func(t *testing.T) {
		return func(t *testing.T) {
			sites := make(map[string]meta.Site)
			for _, s := range sites {
				sites[s.Station+"/"+s.Location] = s
			}
			for _, c := range features {
				if _, ok := sites[c.Station+"/"+c.Location]; !ok {
					t.Error("error: unable to find feature site: " + c.Station + "/" + c.Location)
				}
			}
			for _, c := range features {
				if s, ok := sites[c.Station+"/"+c.Location]; ok {
					switch {
					case c.Start.Before(s.Start):
						t.Log("warning: feature start mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.Start.String(),
							"before",
							s.Start.String(),
						}, " "))
					case s.End.Before(time.Now()) && c.End.After(s.End):
						t.Log("warning: feature end mismatch: " + strings.Join([]string{
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

func TestFeatures(t *testing.T) {

	var features meta.FeatureList
	loadListFile(t, "../network/features.csv", &features)

	for k, fn := range testFeatures {
		t.Run(k, fn(features))
	}
}

func TestFeatures_Stations(t *testing.T) {

	var features meta.FeatureList
	loadListFile(t, "../network/features.csv", &features)

	var stations meta.StationList
	loadListFile(t, "../network/stations.csv", &stations)

	for k, fn := range testFeatures_Stations {
		t.Run(k, fn(features, stations))
	}
}

func TestFeatures_Sites(t *testing.T) {

	var features meta.FeatureList
	loadListFile(t, "../network/features.csv", &features)

	var sites meta.SiteList
	loadListFile(t, "../network/sites.csv", &sites)

	for k, fn := range testFeatures_Sites {
		t.Run(k, fn(features, sites))
	}
}
