package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta/sqlite"
	"github.com/GeoNet/delta/resp"

	_ "modernc.org/sqlite"
)

type Settings struct {
	debug bool // output more operational info

	base string // base directory of delta files on disk
	resp string // base directory of resp files on disk

	db       string // name of the database file
	response string // name of the database response table

	init bool // should the database be updated
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a DELTA Sqlite DB file\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.BoolVar(&settings.debug, "debug", false, "add extra operational info")
	flag.BoolVar(&settings.init, "init", false, "initialise the database if a file on disk")
	flag.StringVar(&settings.base, "base", "", "base directory of delta files on disk")
	flag.StringVar(&settings.resp, "resp", "", "base directory of resp files on disk")
	flag.StringVar(&settings.db, "db", "", "name of the database file on disk")
	flag.StringVar(&settings.response, "response", "Response", "optional database response table name to use")

	flag.Parse()

	ctx := context.Background()

	// set recovers the delta tables
	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatal(err)
	}

	// resp recovers the response files
	files, err := resp.ListBase(settings.resp)
	if err != nil {
		log.Fatal(err)
	}

	values := make(map[string]string)
	for _, f := range files {
		lookup, err := resp.LookupBase(settings.resp, f)
		if err != nil {
			log.Fatal(err)
		}
		values[f] = string(lookup)
	}

	path := ":memory:"
	if settings.db != "" {
		path = settings.db
	}

	opts := url.Values{}
	opts.Set("_time_format", "sqlite")
	opts.Set("_foreign_keys", "on")
	if settings.db != "" && !settings.init {
		opts.Set("mode", "ro")
	}

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?%s", path, url.QueryEscape(opts.Encode())))
	if err != nil {
		log.Fatalf("unable to open database %s: %v", path, err)
	}
	defer db.Close()

	if settings.db == "" || settings.init {

		// insert extra response files
		extra := set.KeyValue(settings.response, "Response", "XML", values)

		log.Println("initialise database")
		start := time.Now()
		if err := sqlite.New(db).Init(ctx, set.TableList(extra)); err != nil {
			log.Fatalf("unable to run database exec: %v", err)
		}
		log.Printf("database initialised in %s", time.Since(start).String())
	}
}
