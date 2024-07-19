package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GeoNet/delta/meta"
)

type Settings struct {
	base   string
	header string
	footer string
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Provide tidal templating\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] [templates ....]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&settings.base, "base", "", "base of delta files on disk")
	flag.StringVar(&settings.header, "header", "", "output header before processing any templates")
	flag.StringVar(&settings.footer, "footer", "", "output footer after all templates have been processed")

	flag.Parse()

	set, err := meta.NewBase(settings.base)
	if err != nil {
		log.Fatal(err)
	}

	if settings.header != "" {
		fmt.Fprintln(os.Stdout, settings.header)
	}

	// process each template given on the command line
	for _, file := range flag.Args() {
		conf, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		config, err := Parse(set, conf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintln(os.Stdout, string(config))
	}

	if settings.footer != "" {
		fmt.Fprintln(os.Stdout, settings.footer)
	}
}
