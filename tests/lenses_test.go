package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestLenses(t *testing.T) {

	var lenses meta.InstalledLensList
	loadListFile(t, "../install/lenses.csv", &lenses)

	t.Run("check for lens duplication", func(t *testing.T) {
		for i := 0; i < len(lenses); i++ {
			for j := i + 1; j < len(lenses); j++ {
				if lenses[i].Mount != lenses[j].Mount {
					continue
				}
				if lenses[i].View != lenses[j].View {
					continue
				}
				if lenses[i].Type != lenses[j].Type {
					continue
				}
				if lenses[i].End.Before(lenses[j].Start) {
					continue
				}
				if lenses[i].Start.After(lenses[j].End) {
					continue
				}
				if lenses[i].End.Equal(lenses[j].Start) {
					continue
				}
				if lenses[i].Start.Equal(lenses[j].End) {
					continue
				}

				t.Errorf("lens duplication: %s/%s (%s)", lenses[i].Mount, lenses[i].View, lenses[i].Start.String())
			}
		}
	})
}
