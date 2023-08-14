package resp

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/GeoNet/delta/internal/stationxml"

	"github.com/google/go-cmp/cmp"
)

func TestAuto(t *testing.T) {
	names, err := fs.Glob(os.DirFS("auto"), "*.xml")
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range names {
		name := name
		t.Run("check resp file: "+name, func(t *testing.T) {
			t.Parallel()

			raw, err := fs.ReadFile(files, filepath.Join("auto", name))
			if err != nil {
				t.Fatal(err)
			}

			first, err := stationxml.NewResponseType(raw)
			if err != nil {
				t.Fatal(err)
			}

			data, err := first.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			second, err := stationxml.NewResponseType(data)
			if err != nil {
				t.Fatal(err)
			}

			check, err := second.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(data, check) {
				t.Error(cmp.Diff(data, check))
			}
		})
	}
}

func TestFiles(t *testing.T) {
	names, err := fs.Glob(os.DirFS("files"), "*.xml")
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range names {
		t.Run("check resp file: "+name, func(t *testing.T) {

			raw, err := fs.ReadFile(files, filepath.Join("files", name))
			if err != nil {
				t.Fatal(err)
			}

			snippet, err := stationxml.NewResponseType(raw)
			if err != nil {
				t.Fatal(err)
			}

			data, err := snippet.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(raw, data) {
				t.Error(cmp.Diff(raw, data))
			}
		})
	}
}
