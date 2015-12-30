package meta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestFilter(t *testing.T) {

	t.Log("Load filter file")
	{
		f, err := meta.LoadFilterFile("../responses/filter.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(f)
	}

}
