package delta

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"sync"

	_ "modernc.org/sqlite"

	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/meta/sqlite"
	"github.com/GeoNet/delta/resp"
)

const ResponseName = "Response"

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

// NewResp returns a slice of response file names.
func NewResp(base string) ([]string, error) {
	return resp.ListBase(base)
}

// DB is a wrapper for an SQL DB pointer.
type DB struct {
	*sql.DB
}

// NewDB returns a DB pointer, a non-empty path is used as a file name, otherwise the database is generated in memory.
func NewDB(path ...string) (*DB, error) {

	file := ":memory:"
	for _, p := range path {
		file = p
	}

	opts := url.Values{}
	opts.Set("_time_format", "sqlite")
	opts.Set("_foreign_keys", "on")

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?%s", file, url.QueryEscape(opts.Encode())))
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

// Init updates the contents of the given database with contents from the given delta base and response paths.
// If these are empty then the default compiled in versions will be used.
func (db *DB) Init(ctx context.Context, set *meta.Set, files ...string) error {

	values := make(map[string]string)
	for _, file := range files {
		lookup, err := resp.Lookup(file)
		if err != nil {
			return err
		}
		values[file] = string(lookup)
	}

	// insert any extra response files
	extra := set.KeyValue(ResponseName, "Response", "XML", values)

	if err := sqlite.New(db.DB).Init(ctx, set.TableList(extra)); err != nil {
		return err
	}

	return nil
}
