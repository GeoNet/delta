package tides

import (
	"reflect"
	"testing"
)

func TestTide(t *testing.T) {
	for _, g := range _tides {
		t.Run("check gauge: "+g.Code, func(t *testing.T) {
			r := Lookup(g.Code)
			if r == nil {
				t.Fatalf("unable to lookup tide: %s", g.Code)
			}
			if !reflect.DeepEqual(&g, r) {
				t.Errorf("unable to match tide: expected %v, found %v", g, r)
			}
		})
	}

	t.Run("check bad gauge", func(t *testing.T) {
		if g := Lookup("XXXX"); g != nil {
			t.Fatal("shouldn't be able to lookup XXXX")
		}
	})
}
