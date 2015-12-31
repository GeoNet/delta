package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/resp"
)

func TestResponse(t *testing.T) {

	t.Log("Load response file")
	{
		r, err := resp.LoadResponseFile("../responses/response.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(r)
	}

}
