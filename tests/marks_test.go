package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var markChecks = map[string]func(*meta.Set) func(t *testing.T){

	// check for session overlaps, there can't be two sessions running at the same mark for the same sampling interval.
	"check for duplicated marks": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			marks := set.Marks()
			for i := 0; i < len(marks); i++ {
				for j := i + 1; j < len(marks); j++ {
					if marks[i].Code == marks[j].Code {
						t.Errorf("mark duplication: " + marks[i].Code)
					}
				}
			}
		}
	},
}

func TestMarks(t *testing.T) {
	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range markChecks {
		t.Run(k, v(set))
	}
}
