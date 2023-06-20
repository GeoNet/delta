package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var noteChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for duplicated notes": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			notes := set.Notes()
			for i := 0; i < len(notes); i++ {
				for j := i + 1; j < len(notes); j++ {
					if notes[i].Code != notes[j].Code {
						continue
					}
					if notes[i].Network != notes[j].Network {
						continue
					}
					if notes[i].Entry != notes[j].Entry {
						continue
					}
					t.Errorf("note duplication: %s/%s", notes[i].Code, notes[i].Network)
				}
			}
		}
	},
}

func TestNotes(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range noteChecks {
		t.Run(k, v(set))
	}
}
