package ntrip

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Formatter interface {
	Decoder
	Encoder
}

func TestFiles_Decode(t *testing.T) {
	files := []struct {
		filename  string
		formatter Formatter
	}{
		{
			FormatsFile,
			&Formats{},
		},
		{
			ModelsFile,
			&Models{},
		},
		{
			UsersFile,
			&Users{},
		},
		{
			MountsFile,
			&Mounts{},
		},
		{
			AliasesFile,
			&Aliases{},
		},
	}

	for _, f := range files {
		t.Run("file decode: "+f.filename, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("testdata", f.filename))
			if err != nil {
				t.Fatal(err)
			}
			if err := ReadBytes(raw, f.formatter); err != nil {
				t.Fatal(err)
			}
			ans, err := WriteBytes(f.formatter)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(raw, ans) {
				t.Errorf("byte mismatch for %s: \"%s\" \"%s\"", f.filename, string(raw), string(ans))
			}
		})
	}
}

func TestFiles_Read(t *testing.T) {
	files := []struct {
		filename  string
		formatter Formatter
	}{
		{
			FormatsFile,
			&Formats{},
		},
		{
			ModelsFile,
			&Models{},
		},
		{
			UsersFile,
			&Users{},
		},
		{
			MountsFile,
			&Mounts{},
		},
		{
			AliasesFile,
			&Aliases{},
		},
	}

	for _, f := range files {
		t.Run("file read: "+f.filename, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("testdata", f.filename))
			if err != nil {
				t.Fatal(err)
			}
			if err := ReadFile(filepath.Join("testdata", f.filename), f.formatter); err != nil {
				t.Fatal(err)
			}
			ans, err := WriteBytes(f.formatter)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(raw, ans) {
				t.Errorf("unexpected %s content -got/+exp\n%s", f.filename, cmp.Diff(string(raw), string(ans)))
			}
		})
	}
}
