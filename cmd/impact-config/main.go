package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"
)

type Settings struct {
	base     string
	resp     string
	channels regexp.Regexp
	skip     regexp.Regexp
	output   string
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build an impact json file from delta meta & response information\n")
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

	flag.StringVar(&settings.base, "base", "", "delta base files")
	flag.StringVar(&settings.resp, "resp", "", "delta resp files")
	flag.TextVar(&settings.channels, "channels", regexp.MustCompile("^[EBH][NH]Z$"), "match impact channels")
	flag.TextVar(&settings.skip, "skip", regexp.MustCompile("^SB$"), "skip networks")
	flag.StringVar(&settings.output, "output", "", "output impact json file")

	flag.Parse()

	set, err := meta.NewBase(settings.base)
	if err != nil {
		log.Fatal(err)
	}

	streams, err := settings.ImpactStreams(set, resp.NewResp(settings.resp))
	if err != nil {
		log.Fatalf("problem loading streams %s: %v\n", settings.base, err)
	}

	res, err := json.MarshalIndent(streams, "", "  ")
	if err != nil {
		log.Fatalf("problem marshalling streams %s: %v\n", settings.base, err)
	}

	switch settings.output {
	case "", "-":
		os.Stdout.Write(res)
	default:
		if err := os.WriteFile(settings.output, res, 0600); err != nil {
			log.Fatalf("error: unable to write file %s: %v", settings.output, err)
		}
	}
}
