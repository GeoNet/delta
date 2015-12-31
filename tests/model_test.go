package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/resp"
)

func TestModel(t *testing.T) {

	t.Log("Load model file")
	{
		m, err := resp.LoadModelFile("../responses/model.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(m)
	}

}
