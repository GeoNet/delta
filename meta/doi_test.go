package meta

import (
	"testing"
)

func TestDOI(t *testing.T) {

	good := []string{
		"http://doi.org/10.21420/8TCZ-TV02",
		"http://doi.org/10.1029/2020EO144274",
		"https://doi.org/10.21420/8TCZ-TV02",
		"https://doi.org/10.1029/2020EO144274",
		"doi.org/10.21420/8TCZ-TV02",
		"doi.org/10.1029/2020EO144274",
		"10.21420/8TCZ-TV02",
		"10.1029/2020EO144274",
		"DOI:10.21420/8TCZ-TV02",
		"DOI:10.1029/2020EO144274",
	}

	for _, v := range good {
		d, err := NewDoi(v)
		if err != nil {
			t.Errorf("unable to pass doi %s: %v", v, err)
		}
		if q := d.String(); q != v {
			t.Errorf("unexpected doi check, expected %q, got %q", v, q)

		}
	}

	bad := []string{
		"",
		"//doi.org/10.21420/8TCZ-TV02",
		"http://xdoi.org/10.1029/2020EO144274",
		"xdoi.org/10.21420/8TCZ-TV02",
		"doi.org /10.1029/2020EO144274",
		"10.21420/ 8TCZ-TV02",
		"1029/2020EO144274",
		"8TCZ-TV02",
	}

	for _, v := range bad {
		if _, err := NewDoi(v); err == nil {
			t.Errorf("able to pass bad doi %s", v)
		}
	}

}
