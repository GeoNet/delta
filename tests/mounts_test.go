package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestMounts(t *testing.T) {

	var mounts meta.MountList
	loadListFile(t, "../network/mounts.csv", &mounts)

	t.Run("check for mount duplication", func(t *testing.T) {
		for i := 0; i < len(mounts); i++ {
			for j := i + 1; j < len(mounts); j++ {
				if mounts[i].Code == mounts[j].Code {
					t.Errorf("mount duplication: " + mounts[i].Code)
				}
			}
		}
	})
}
