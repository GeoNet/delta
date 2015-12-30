package meta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestResponse(t *testing.T) {

	t.Log("Load response file")
	{
		r, err := meta.LoadResponseFile("../responses/response.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(r)
	}

}
