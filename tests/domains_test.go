package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var domainChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicate domains": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			domains := make(map[string]interface{})
			for _, s := range set.Domains() {
				if _, ok := domains[s.Name]; ok {
					t.Errorf("domain %s is duplicated", s.Name)
				}
				domains[s.Name] = true
			}
		}
	},
}

func TestDomains(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range domainChecks {
		t.Run(k, v(set))
	}
}
