package meta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestFIR(t *testing.T) {

	t.Log("Load fir file")
	{
		f, err := meta.LoadFIRFile("../responses/fir.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(f)
	}

}
