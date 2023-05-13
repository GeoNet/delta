package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GeoNet/delta"
)

type Settings struct {
	baseDir    string
	outputFile string

	dartNetworks    string
	coastalNetworks string
	enviroNetworks  string
	manualNetworks  string
}

func (s Settings) split(str string) []string {
	var list []string
	for _, s := range strings.Split(str, ",") {
		if s := strings.TrimSpace(s); s != "" {
			list = append(list, s)
		}
	}
	return list
}

func (s Settings) Dart() []string    { return s.split(s.dartNetworks) }
func (s Settings) Coastal() []string { return s.split(s.coastalNetworks) }
func (s Settings) Enviro() []string  { return s.split(s.enviroNetworks) }
func (s Settings) Manual() []string  { return s.split(s.manualNetworks) }

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a tilde domain config file\n")
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

	// application settings
	flag.StringVar(&settings.baseDir, "base", "", "delta base files")
	flag.StringVar(&settings.outputFile, "output", "", "output tilde configuration file")

	flag.StringVar(&settings.dartNetworks, "dart", "TD", "dart buoy network code")
	flag.StringVar(&settings.coastalNetworks, "coastal", "TG,LG", "coast tsunami gauge network code")
	flag.StringVar(&settings.enviroNetworks, "enviro", "EN", "envirosensor network code")
	flag.StringVar(&settings.manualNetworks, "manual", "MC", "manualcollect network code")

	flag.Parse()

	// set recovers the delta tables
	set, err := delta.NewBase(settings.baseDir)
	if err != nil {
		log.Fatal(err)
	}

	var tilde Tilde

	// update coastal domain
	if err := tilde.Coastal(set, settings.Coastal()...); err != nil {
		log.Fatal(err)
	}

	// update dart domain
	if err := tilde.Dart(set, settings.Dart()...); err != nil {
		log.Fatal(err)
	}

	// update envirosensor domain
	if err := tilde.EnviroSensor(set, settings.Enviro()...); err != nil {
		log.Fatal(err)
	}

	// update manualcollection domain
	if err := tilde.ManualCollection(set, settings.Manual()...); err != nil {
		log.Fatal(err)
	}

	// keep things in order
	tilde.Sort()

	switch {
	case settings.outputFile != "":
		// output file has been given
		file, err := os.Create(settings.outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		if err := tilde.MarshalIndent(file, "", "  "); err != nil {
			log.Fatal(err)
		}
	default:
		if err := tilde.MarshalIndent(os.Stdout, "", "  "); err != nil {
			log.Fatal(err)
		}
	}
}
