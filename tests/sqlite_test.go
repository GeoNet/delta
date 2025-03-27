package delta_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta/sqlite"
	"github.com/GeoNet/delta/resp"
)

func TestSqLite(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// resp recovers the response files
	files, err := resp.List()
	if err != nil {
		t.Fatal(err)
	}

	values := make(map[string]string)
	for _, f := range files {
		lookup, err := resp.Lookup(f)
		if err != nil {
			t.Fatal(err)
		}
		values[f] = string(lookup)
	}

	path := ":memory:"

	opts := url.Values{}
	opts.Set("_time_format", "sqlite")
	opts.Set("_foreign_keys", "on")

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?%s", path, url.QueryEscape(opts.Encode())))
	if err != nil {
		t.Fatalf("unable to open database %s: %v", path, err)
	}
	defer db.Close()

	// insert extra response files
	extra := set.KeyValue("Response", "Response", "XML", values)

	t.Log("initialise database")
	start := time.Now()
	if err := sqlite.New(db).Init(ctx, set.TableList(extra)); err != nil {
		t.Fatalf("unable to run database exec: %v", err)
	}
	t.Logf("database initialised in %s", time.Since(start).String())

}
