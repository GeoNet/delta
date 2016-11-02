package resp

import (
	"testing"
)

func TestResp_StageSet(t *testing.T) {

	var tests = []struct {
		t string
		s StageSet
	}{
		{"a2d", A2D{}},
		{"paz", PAZ{}},
		{"fir", FIR{}},
		{"poly", Polynomial{}},
	}

	for _, test := range tests {
		if test.t != test.s.GetType() {
			t.Errorf("invalid stage type: got %v, expected %v", test.t, test.s.GetType())
		}
	}
}
