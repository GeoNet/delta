package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var citationChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicate citations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			citations := make(map[string]interface{})
			for _, c := range set.Citations() {
				if _, ok := citations[c.Key]; ok {
					t.Errorf("citation %s is duplicated", c.Key)
				}
				citations[c.Key] = true
			}
		}
	},
}

func TestCitations(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range citationChecks {
		t.Run(k, v(set))
	}
}
