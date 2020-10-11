package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testMarks = map[string]func([]meta.Mark) func(t *testing.T){

	// check for session overlaps, there can't be two sessions running at the same mark for the same sampling interval.
	"check for duplicated marks": func(marks []meta.Mark) func(t *testing.T) {
		return func(t *testing.T) {
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
	var marks meta.MarkList
	loadListFile(t, "../network/marks.csv", &marks)

	for k, fn := range testMarks {
		t.Run(k, fn(marks))
	}
}
