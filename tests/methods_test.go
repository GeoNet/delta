package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var methodChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicate methods": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			methods := make(map[meta.Method]interface{})
			for _, c := range set.Methods() {
				key := meta.Method{
					Domain: c.Domain,
					Name:   c.Name,
				}
				if _, ok := methods[key]; ok {
					t.Errorf("method %s/%s is duplicated", c.Domain, c.Name)
				}
				methods[key] = true
			}
		}
	},
}

func TestMethods(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range methodChecks {
		t.Run(k, v(set))
	}
}
