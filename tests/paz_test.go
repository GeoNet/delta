package meta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestPAZ(t *testing.T) {

	t.Log("Load fir file")
	{
		pz, err := meta.LoadPAZFile("../responses/paz.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(pz)
	}

}
