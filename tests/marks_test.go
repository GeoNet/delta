package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestMarks(t *testing.T) {

	var marks meta.MarkList
	t.Log("Load network marks file")
	{
		if err := meta.LoadList("../network/marks.csv", &marks); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < len(marks); i++ {
		for j := i + 1; j < len(marks); j++ {
			if marks[i].Code == marks[j].Code {
				t.Errorf("mark duplication: " + marks[i].Code)
			}
		}
	}
}
