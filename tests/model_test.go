package meta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestModel(t *testing.T) {

	t.Log("Load model file")
	{
		m, err := meta.LoadModelFile("../responses/model.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(m)
	}

}
