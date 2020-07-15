package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestViews(t *testing.T) {

	var views meta.ViewList
	loadListFile(t, "../network/views.csv", &views)

	t.Run("check for view duplication", func(t *testing.T) {
		for i := 0; i < len(views); i++ {
			for j := i + 1; j < len(views); j++ {
				if views[i].Mount != views[j].Mount {
					continue
				}
				if views[i].Code == views[j].Code {
					t.Errorf("view duplication: %s/%s", views[i].Mount, views[i].Code)
				}
			}
		}
	})
}
