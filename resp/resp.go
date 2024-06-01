package resp

import (
	"embed"
	"encoding/xml"
	"fmt"
	"io/fs"
	"os"
	"sync"
)

// TODO: add embed when populated
var locations = []string{"files", "auto"}

//go:embed files/*.xml
//go:embed auto/*.xml
var files embed.FS

// LookupFS returns a byte slice representation of a generic stationxml Response if present in the given file system.
func LookupFS(fsys fs.FS, response string) ([]byte, error) {
	for _, l := range locations {
		names, err := fs.Glob(fsys, fmt.Sprintf("%s/%s.xml", l, response))
		if err != nil {
			return nil, err
		}

		// return the first one found
		for _, name := range names {
			data, err := fs.ReadFile(fsys, name)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}

	return nil, nil
}

// Lookup returns a byte slice representation of a generic embeded stationxml Response if present.
func Lookup(response string) ([]byte, error) {
	return LookupFS(files, response)
}

// LookupDir returns a byte slice representation of a generic stationxml Response if stored in the given directory.
func LookupDir(path string, response string) ([]byte, error) {
	return LookupFS(os.DirFS(path), response)
}

// LookupBase returns a byte slice representation of a generic stationxml Response either stored in a given directory or
// in the embedded files if no base directory given.
func LookupBase(base string, response string) ([]byte, error) {
	if base != "" {
		return LookupDir(base, response)
	}
	return Lookup(response)
}

type Resp struct {
	mu sync.Mutex

	base  string
	cache map[string][]byte
}

func NewResp(base string) *Resp {
	return &Resp{
		base:  base,
		cache: make(map[string][]byte),
	}
}

func (r *Resp) Lookup(response string) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if v, ok := r.cache[response]; ok {
		return v, nil
	}

	v, err := LookupBase(r.base, response)
	if err != nil {
		return nil, err
	}

	r.cache[response] = v

	return v, nil
}

func (r *Resp) Type(response string) (*ResponseType, error) {

	data, err := r.Lookup(response)
	if err != nil {
		return nil, err
	}

	// decode the response into a simple form.
	var res ResponseType
	if err := xml.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
