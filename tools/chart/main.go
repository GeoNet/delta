package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {

	var base string
	flag.StringVar(&base, "base", "../..", "delta base files")

	var output string
	flag.StringVar(&output, "output", ".", "output xml directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build an altus XML file from delta meta & response information\n")
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

	flag.Parse()

	for _, path := range flag.Args() {

		charts, err := ReadPlots(path)
		if err != nil {
			log.Fatalf("problem loading config file %s: %v", path, err)
		}
		for plot, config := range charts.Configs {
			var pages []Page
			var plots []Plot
			for _, p := range config.Pages {
				switch p.Type {
				case "tsunami":
					res, err := p.Tsunami(base)
					if err != nil {
						log.Fatalf("problem build tsunami pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "drum":
					res, err := p.Drum(base)
					if err != nil {
						log.Fatalf("problem build drum pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "trace":
					res, err := p.Trace(base)
					if err != nil {
						log.Fatalf("problem build trace pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "depth":
					res, err := p.Depth(base)
					if err != nil {
						log.Fatalf("problem build depth pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "combined":
					res, err := p.Combined(base)
					if err != nil {
						log.Fatalf("problem build combined pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "combined-medium":
					res, err := p.CombinedMedium(base)
					if err != nil {
						log.Fatalf("problem build combined-medium pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "traces":
					res, err := p.Traces(base)
					if err != nil {
						log.Fatalf("problem build traces pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "drum-small":
					res, err := p.DrumSmall(base)
					if err != nil {
						log.Fatalf("problem build drum-small pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "gauge":
					res, err := p.Gauge(base)
					if err != nil {
						log.Fatalf("problem build gauge pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "temperature":
					res, err := p.Temperature(base)
					if err != nil {
						log.Fatalf("problem build temperature pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "networks":
					res, err := p.Networks(base)
					if err != nil {
						log.Fatalf("problem build networks pages %s: %v", plot, err)
					}
					pages = append(pages, res...)
				case "wave":
					res, err := p.Wave(base)
					if err != nil {
						log.Fatalf("problem build wave pages %s: %v", plot, err)
					}
					plots = append(plots, res...)
				default:
					log.Fatalf("unknown plot type %s", p.Type)
				}
			}

			settings := Chart{
				Pages: pages,
				Plots: plots,
			}

			res, err := settings.Marshal()
			if err != nil {
				log.Fatalf("error: unable to marshal xml: %v", err)
			}

			outfile := filepath.Join(output, config.Filename)
			if err := os.MkdirAll(filepath.Dir(outfile), 0755); err != nil {
				log.Fatalf("error: unable to create directory %s: %v", filepath.Dir(output), err)
			}
			if err := ioutil.WriteFile(outfile, res, 0644); err != nil {
				log.Fatalf("error: unable to write file %s: %v", outfile, err)
			}
		}
	}
}
