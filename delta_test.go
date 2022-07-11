package delta

import (
	"testing"
)

func TestDelta(t *testing.T) {
	if _, err := New(); err != nil {
		t.Fatal(err)
	}
}
