package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/resp"
)

func TestFIR(t *testing.T) {

	t.Log("Load fir file")
	{
		f, err := resp.LoadFIRFile("../responses/fir.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(f)
	}

}
