package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestMarks(t *testing.T) {
	var marks meta.MarkList
	loadListFile(t, "../network/marks.csv", &marks)

	t.Run("check for duplicated marks", func(t *testing.T) {
		for i := 0; i < len(marks); i++ {
			for j := i + 1; j < len(marks); j++ {
				if marks[i].Code == marks[j].Code {
					t.Errorf("mark duplication: " + marks[i].Code)
				}
			}
		}
	})
}
