package tides

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTide(t *testing.T) {
	for _, g := range _tides {
		t.Run("check gauge: "+g.Code, func(t *testing.T) {
			r := Lookup(g.Code)
			if r == nil {
				t.Fatalf("unable to lookup tide: %s", g.Code)
			}
			if !cmp.Equal(g, *r) {
				t.Errorf("unable to match tide %s: %s", g.Code, cmp.Diff(g, *r))
			}
		})
	}

	t.Run("check bad gauge", func(t *testing.T) {
		if g := Lookup("XXXX"); g != nil {
			t.Fatal("shouldn't be able to lookup XXXX")
		}
	})
}
