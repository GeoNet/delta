package meta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestPolynomial(t *testing.T) {

	t.Log("Load polynomial file")
	{
		p, err := meta.LoadPolynomialFile("../responses/polynomial.toml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(p)
	}

}
