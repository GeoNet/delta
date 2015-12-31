package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/resp"
)

func TestPAZ(t *testing.T) {

	t.Log("Load fir file")
	{
		pz, err := resp.LoadPAZFile("../responses/paz.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(pz)
	}

}
