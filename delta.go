package delta

import (
	"embed"
	"io/fs"
	"os"
	"sync"

	"github.com/GeoNet/delta/meta"
)

//go:embed assets/*.csv
//go:embed install/*.csv
//go:embed network/*.csv
//go:embed environment/*.csv
//go:embed references/*.csv
var files embed.FS

// there can be but one
var singleton struct {
	once sync.Once
	set  *meta.Set
	err  error
}

// New returns a Set pointer based on the embeded csv files.
// Multiple calls will return the same pointer and error as the first call.
func New() (*meta.Set, error) {
	singleton.once.Do(func() {
		singleton.set, singleton.err = NewFS(files)
	})
	return singleton.set, singleton.err
}

// NewFS returns a Delta pointer based on a given FS structure.
func NewFS(fs fs.FS) (*meta.Set, error) {

	set, err := meta.NewSet(fs)
	if err != nil {
		return nil, err
	}

	return set, nil
}

// NewBase returns a Delta pointer based on an optional directory base prefix.
// If the base is empty then the default embeded Set will be returned.
func NewBase(base string) (*meta.Set, error) {
	if base != "" {
		return NewFS(os.DirFS(base))
	}
	return New()
}
