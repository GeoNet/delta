package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta/sqlite"

	_ "modernc.org/sqlite"
)

type Settings struct {
	debug bool // output more operational info

	base string // base directory of delta files on disk

	db     string // name of the database file
	schema string // name of the database schema

	init     bool   // should the database be updated
	listen   bool   // should a web service be enabled
	hostport string // hostport to listen on for the web service
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a DELTA Sqlite DB and optional REST API service\n")
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
	flag.BoolVar(&settings.listen, "listen", false, "should a web service be enabled")

	flag.StringVar(&settings.base, "base", "", "base directory of delta files on disk")
	flag.StringVar(&settings.db, "db", "", "name of the database file on disk")
	flag.StringVar(&settings.schema, "schema", "", "optional database schema to use")
	flag.StringVar(&settings.hostport, "hostport", ":8080", "base directory of delta files on disk")

	flag.Parse()

	ctx := context.Background()

	// set recovers the delta tables
	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatal(err)
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
		log.Println("initialise database")
		start := time.Now()
		if err := sqlite.New(db, settings.schema).Init(ctx, set.TableList()); err != nil {
			log.Fatalf("unable to run database exec: %v", err)
		}
		log.Printf("database initialised in %s", time.Since(start).String())
	}

	if !settings.listen {
		return
	}

	log.Printf("handling requests on %s", settings.hostport)
	server := &http.Server{
		Addr:              settings.hostport,
		Handler:           newHandler(db),
		ReadHeaderTimeout: 3 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
