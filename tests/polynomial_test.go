package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/resp"
)

func TestPolynomial(t *testing.T) {

	t.Log("Load polynomial file")
	{
		p, err := resp.LoadPolynomialFile("../responses/polynomial.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(p)
	}

}
