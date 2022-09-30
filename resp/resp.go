package resp

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
)

//TODO: add embed when populated
var locations = []string{"files", "auto"}

//go:embed files/*.xml
//go:embed auto/*.xml
var files embed.FS

// LookupFS returns a byte slice of a stationxml Response if present in the given file system.
func LookupFS(fsys fs.FS, response string) ([]byte, error) {
	for _, l := range locations {
		names, err := fs.Glob(fsys, fmt.Sprintf("%s/%s.xml", l, response))
		if err != nil {
			return nil, err
		}

		for _, name := range names {
			data, err := fs.ReadFile(files, name)
			if err != nil {
				return nil, err
			}

			return data, nil
		}
	}

	return nil, nil
}

// Lookup returns a pointer to an embeded stationxml Response if present.
func Lookup(response string) ([]byte, error) {
	return LookupFS(files, response)
}

// LookupDir returns a pointer to an stationxml Response if stored in the given directory.
func LookupDir(path string, response string) ([]byte, error) {
	return LookupFS(os.DirFS(path), response)
}

// LookupBase returns a pointer to an stationxml Response either stored in a given directory or
// in the embedded files if no base directory given.
func LookupBase(base string, response string) ([]byte, error) {
	if base != "" {
		return LookupDir(base, response)
	}
	return Lookup(response)
}
