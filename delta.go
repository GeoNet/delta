package delta

import (
	"embed"
	"io/fs"
	"sync"

	"github.com/GeoNet/delta/meta"
)

//go:embed assets/*.csv
//go:embed install/*.csv
//go:embed network/*.csv
//go:embed environment/*.csv
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
