package meta

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSet(t *testing.T) {
	_, err := NewSet(os.DirFS("testdata"), func(s string) string {
		switch {
		case filepath.Dir(s) == "assets":
			return "assets.csv"
		default:
			return filepath.Base(s)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
}
