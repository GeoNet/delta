package meta

import (
	"testing"
)

func TestParseInt(t *testing.T) {

	good := map[string]int{
		"":    0,
		" ":   0,
		"0":   0,
		" 0 ": 0,
		"10":  10,
		"-10": -10,
	}

	bad := []string{
		"0 0",
		"9.9",
		"A",
	}

	for k, v := range good {
		switch n, err := ParseInt(k); {
		case err != nil:
			t.Fatal(err)
		case n != v:
			t.Errorf("problem parsing int \"%s\", expected %d got %d", k, v, n)
		}
	}

	for _, k := range bad {
		if _, err := ParseInt(k); err == nil {
			t.Errorf("problem not parsing int \"%s\", expected an error", k)
		}
	}

}

func TestParseFloat64(t *testing.T) {

	good := map[string]float64{
		"":        0.0,
		" ":       0.0,
		" 1.0 ":   1.0,
		"1.0":     1.0,
		"-1.0":    -1.0,
		"1.0e+01": 10.0,
	}

	bad := []string{
		"0 0",
		"A",
		"-NaN",
	}

	for k, v := range good {
		switch n, err := ParseFloat64(k); {
		case err != nil:
			t.Error(err)
		case n != v:
			t.Errorf("problem parsing float \"%s\", expected %g got %g", k, v, n)
		}
	}

	for _, k := range bad {
		if _, err := ParseFloat64(k); err == nil {
			t.Errorf("problem not parsing float \"%s\", expected an error", k)
		}
	}

}

func TestParseSamplingRate(t *testing.T) {

	good := map[string]float64{
		"":         0.0,
		" ":        0.0,
		"1":        1.0,
		"1.0e+01":  10.0,
		"-1":       1.0,
		"-10":      1.0 / 10.0,
		"-1.0e+01": 1.0 / 10.0,
		"-1.0e-01": 10.0,
	}

	bad := []string{
		"0 0",
		"A",
		"-NaN",
	}

	for k, v := range good {
		switch n, err := ParseSamplingRate(k); {
		case err != nil:
			t.Fatal(err)
		case n != v:
			t.Errorf("problem parsing float \"%s\", expected %g got %g", k, v, n)
		}
	}

	for _, k := range bad {
		if _, err := ParseSamplingRate(k); err == nil {
			t.Errorf("problem not parsing float \"%s\", expected an error", k)
		}
	}

}
