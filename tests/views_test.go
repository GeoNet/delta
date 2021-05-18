package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testViews = map[string]func([]meta.View) func(t *testing.T){

	"check for duplicated views": func(views []meta.View) func(t *testing.T) {
		return func(t *testing.T) {
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
		}
	},
}

func TestViews(t *testing.T) {

	var views meta.ViewList
	loadListFile(t, "../network/views.csv", &views)

	for k, fn := range testViews {
		t.Run(k, fn(views))
	}
}
