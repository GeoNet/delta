package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/resp"
)

func TestModel(t *testing.T) {

	t.Log("Load sensor model file")
	{
		m, err := resp.LoadSensorModelFile("../responses/sensor.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(m)
	}

}
