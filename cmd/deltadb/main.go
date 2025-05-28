package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GeoNet/delta"
)

type Settings struct {
	base string // base directory of delta files on disk
	resp string // base directory of resp files on disk
	path string // name of the database file
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build and initialise a DELTA Sqlite database\n")
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

	flag.StringVar(&settings.base, "base", "", "base directory of delta files on disk, default uses embedded files")
	flag.StringVar(&settings.resp, "resp", "", "base directory of resp files on disk, default uses embedded files")
	flag.StringVar(&settings.path, "path", "", "name of the database file on disk, default is to use memory only")

	flag.Parse()

	ctx := context.Background()

	// set recovers the delta tables
	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatal(err)
	}

	// recover any response files
	files, err := delta.NewResp(settings.resp)
	if err != nil {
		log.Fatal(err)
	}

	// open the database file handle
	db, err := delta.NewDB(settings.path)
	if err != nil {
		log.Fatalf("unable to open database %q: %v", settings.path, err)
	}
	defer db.Close()

	start := time.Now()

	// initialise the database
	if err := db.Init(ctx, set, files...); err != nil {
		log.Fatalf("unable to init database: %v", err)
	}

	switch {
	case settings.path != "":
		log.Printf("successfully initialised %q in %s", settings.path, time.Since(start).Truncate(time.Millisecond).String())
	default:
		log.Printf("successfully initialised memory in %s", time.Since(start).Truncate(time.Millisecond).String())
	}
}
