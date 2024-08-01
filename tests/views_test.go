package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var viewChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for duplicated views": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			views := set.Views()
			for i := 0; i < len(views); i++ {
				for j := i + 1; j < len(views); j++ {
					if views[i].Mount != views[j].Mount {
						continue
					}
					if views[i].Code != views[j].Code {
						continue
					}
					if views[i].Start.After(views[j].End) {
						continue
					}
					if views[i].End.Before(views[j].Start) {
						continue
					}
					t.Errorf("view duplication: %s/%s", views[i].Mount, views[i].Code)
				}
			}
		}
	},
}

func TestViews(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range viewChecks {
		t.Run(k, v(set))
	}
}
