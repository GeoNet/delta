package expr

import (
	"testing"
)

func TestExpr_Float64(t *testing.T) {

	good := map[string]float64{
		"1.0 + 1.0":       2.0,
		"(1.0 + 1.0)/0.5": 4.0,
	}

	for k, v := range good {
		switch x, err := ToFloat64(k); {
		case err != nil:
			t.Errorf("error with expr \"%s\": %v", k, err)
		case x != v:
			t.Errorf("invalid expr for \"%s\": expected %g, but got %g", k, v, x)
		}
	}

	bad := []string{"", "a", "1.0 / 0.0", "1.0 + ()", "2001.169"}

	for _, v := range bad {
		if _, err := ToFloat64(v); err == nil {
			t.Errorf("invalid expr for \"%s\": expected and error but was nil", v)
		}
	}
}

func TestExpr_Int64(t *testing.T) {

	good := map[string]int64{
		"1":             1,
		"-1":            -1,
		"- 1":           -1,
		"1 + 1":         2,
		"(1 + 3)*2 + 2": 10,
		"43 / 2":        21,
	}

	for k, v := range good {
		switch x, err := ToInt64(k); {
		case err != nil:
			t.Errorf("error with expr \"%s\": %v", k, err)
		case x != v:
			t.Errorf("invalid expr for \"%s\": expected %d, but got %d", k, v, x)
		}

		switch x, err := ToInt(k); {
		case err != nil:
			t.Errorf("error with expr \"%s\": %v", k, err)
		case int64(x) != v:
			t.Errorf("invalid expr for \"%s\": expected %d, but got %d", k, v, x)
		}
	}

	bad := []string{"", " ", "\t", "1/0", "1 +"}

	for _, v := range bad {
		if _, err := ToInt64(v); err == nil {
			t.Errorf("invalid expr for \"%s\": expected and error but was nil", v)
		}
	}
}

func TestExpr_Int(t *testing.T) {

	good := map[string]int{
		"1":             1,
		"-1":            -1,
		"- 1":           -1,
		"1 + 1":         2,
		"(1 + 3)*2 + 2": 10,
		"43 / 2":        21,
	}

	for k, v := range good {
		switch x, err := ToInt(k); {
		case err != nil:
			t.Errorf("error with expr \"%s\": %v", k, err)
		case x != v:
			t.Errorf("invalid expr for \"%s\": expected %d, but got %d", k, v, x)
		}
	}

	bad := []string{"", " ", "\t", "1/0", "1 +"}

	for _, v := range bad {
		if _, err := ToInt(v); err == nil {
			t.Errorf("invalid expr for \"%s\": expected and error but was nil", v)
		}
	}
}

func TestExpr_Uint(t *testing.T) {

	good := map[string]uint{
		"1":             1,
		"1 + 1":         2,
		"(1 + 3)*2 + 2": 10,
		"43 / 2":        21,
	}

	for k, v := range good {
		switch x, err := ToUint(k); {
		case err != nil:
			t.Errorf("error with expr \"%s\": %v", k, err)
		case x != v:
			t.Errorf("invalid expr for \"%s\": expected %d, but got %d", k, v, x)
		}
	}

	bad := []string{"", " ", "\t", "1/0", "1 +", "-1"}

	for _, v := range bad {
		if _, err := ToUint(v); err == nil {
			t.Errorf("invalid expr for \"%s\": expected and error but was nil", v)
		}
	}
}

func TestExpr_Bool(t *testing.T) {

	good := map[string]bool{
		"1 == 1":       true,
		"(2 + 2) == 4": true,
		"1 == 2":       false,
	}

	for k, v := range good {
		switch x, err := ToBool(k); {
		case err != nil:
			t.Errorf("error with expr \"%s\": %v", k, err)
		case x != v:
			t.Errorf("invalid expr for \"%s\": expected %v, but got %v", k, v, x)
		}
	}

	bad := []string{"", " ", "\t", "1", "1 + 2", "fake"}

	for _, v := range bad {
		if _, err := ToBool(v); err == nil {
			t.Errorf("invalid expr for \"%s\": expected and error but was nil", v)
		}
	}
}

func TestExpr_String(t *testing.T) {

	good := map[string]string{
		`"s"`:       "s",
		"`s`":       "s",
		`"s" + "s"`: "ss",
	}

	for k, v := range good {
		switch x, err := ToString(k); {
		case err != nil:
			t.Errorf("error with expr \"%s\": %v", k, err)
		case x != v:
			t.Errorf("invalid expr for \"%s\": expected %v, but got %v", k, v, x)
		}
	}

	bad := []string{"", "s + s"}

	for _, v := range bad {
		if _, err := ToString(v); err == nil {
			t.Errorf("invalid expr for \"%s\": expected and error but was nil", v)
		}
	}
}
