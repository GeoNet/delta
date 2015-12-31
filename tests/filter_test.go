package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/resp"
)

func TestFilter(t *testing.T) {

	t.Log("Load filter file")
	{
		f, err := resp.LoadFilterFile("../responses/filter.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(f)
	}

}
