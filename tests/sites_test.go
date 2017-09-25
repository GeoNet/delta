package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestSites(t *testing.T) {

	var sites meta.SiteList
	loadListFile(t, "../network/sites.csv", &sites)

	t.Run("check for duplicated sites", func(t *testing.T) {
		for i := 0; i < len(sites); i++ {
			for j := i + 1; j < len(sites); j++ {
				if sites[i].Station == sites[j].Station && sites[i].Location == sites[j].Location {
					t.Errorf("site duplication: " + sites[i].Station + "/" + sites[i].Location)
				}
			}
		}
	})

}
