package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/GeoNet/delta"
)

type Settings struct {
	base string // options delta base file directory

	networks string // overall installed sensor network code
	gnss     string // GNSS network codes
	dart     string // DART network code
	lentic   string // lake network code
	coastal  string // coastal network code
	magnetic string // magnetic network code
	enviro   string // envirosensor network code
	manual   string // manualcollect network code
	volcano  string // volcano camera network code
	building string // building camera network code
	doas     string // doas network code

	seismic  regexp.Regexp // seismic location codes
	strong   regexp.Regexp // seismic location codes
	acoustic regexp.Regexp // seismic location codes
	combined regexp.Regexp // installed location codes
	water    regexp.Regexp // water location codes
	geomag   regexp.Regexp // geomag location codes

	output string // optional output file
	groups bool   // use groups output format
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a sensor description file\n")
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
	flag.StringVar(&settings.networks, "networks", "AK,CB,CH,EC,FI,HB,KI,NM,NZ,OT,RA,RT,SC,SI,SM,SP,TP,TR,WL", "installed network codes")
	flag.StringVar(&settings.coastal, "coastal", "TG", "coastal tsunami gauge network code")
	flag.StringVar(&settings.lentic, "lentic", "LG", "lentic tsunami gauge network code")
	flag.StringVar(&settings.dart, "dart", "TD", "dart buoy network code")
	flag.StringVar(&settings.enviro, "enviro", "EN", "envirosensor network code")
	flag.StringVar(&settings.manual, "manual", "MC", "manualcollect network code")
	flag.StringVar(&settings.volcano, "volcano", "VC", "volcano camera network code")
	flag.StringVar(&settings.building, "building", "BC", "building camera network code")
	flag.StringVar(&settings.magnetic, "magnetic", "GM,SM", "geomagnetic network code")
	flag.StringVar(&settings.doas, "doas", "EN", "doas network code")
	flag.TextVar(&settings.combined, "combined", regexp.MustCompile("^[123]"), "combined sensor location codes")
	flag.StringVar(&settings.gnss, "gnss", "CG,GN,IG,LI,SA", "GNSS network codes")
	flag.TextVar(&settings.seismic, "seismic", regexp.MustCompile("^1"), "combined sensor location codes")
	flag.TextVar(&settings.strong, "strong", regexp.MustCompile("^2"), "combined sensor location codes")
	flag.TextVar(&settings.acoustic, "acoustic", regexp.MustCompile("^3"), "combined sensor location codes")
	flag.TextVar(&settings.water, "water", regexp.MustCompile("^4"), "water pressue sensor codes")
	flag.TextVar(&settings.geomag, "geomag", regexp.MustCompile("^5"), "geomag sensor codes")

	flag.BoolVar(&settings.groups, "groups", false, "use groups in output XML format")
	flag.StringVar(&settings.output, "output", "", "output sensor description file")

	flag.Parse()

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to load delta base files: %v", err)
	}

	var network Network

	switch {
	case settings.groups:
		network = settings.Groups(set)
	default:
		// this is the legacy mechanism - which has trouble with sensors at sites which have unrelated networks
		for _, n := range strings.Split(settings.networks, ",") {
			if err := network.InstalledSensors(set, &settings.combined, strings.TrimSpace(n), ""); err != nil {
				log.Fatalf("unable to build seismic details (%s): %v", n, err)
			}
		}

		if err := network.InstalledSensors(set, &settings.water, settings.coastal, "Coastal"); err != nil {
			log.Fatalf("unable to build water details: %v", err)
		}

		if err := network.InstalledSensors(set, &settings.water, settings.lentic, "Lake"); err != nil {
			log.Fatalf("unable to build water details: %v", err)
		}

		if err := network.EnviroSensor(set, settings.enviro, "Environmental Sensor"); err != nil {
			log.Fatalf("unable to build envirosensor configuration: %v", err)
		}

		if err := network.Dart(set, settings.dart, "DART Bottom Pressure Recorder"); err != nil {
			log.Fatalf("unable to build dart configuration: %v", err)
		}

		if err := network.Gnss(set, "GNSS Antenna", "GNSS Receiver"); err != nil {
			log.Fatalf("unable to build gnss configuration: %v", err)
		}

		if err := network.Camera(set, settings.volcano, "Camera"); err != nil {
			log.Fatalf("unable to build camera configuration: %v", err)
		}

		if err := network.Camera(set, settings.building, "Camera"); err != nil {
			log.Fatalf("unable to build camera configuration: %v", err)
		}

		if err := network.Doas(set, settings.doas, "DOAS"); err != nil {
			log.Fatalf("unable to build camera configuration: %v", err)
		}

		if err := network.ManualCollection(set, settings.manual, "Manual Collection"); err != nil {
			log.Fatalf("unable to build camera configuration: %v", err)
		}
	}

	switch {
	case settings.output != "":
		// output file has been given
		file, err := os.Create(settings.output)
		if err != nil {
			log.Fatalf("unable to create output file %q: %v", settings.output, err)
		}
		defer file.Close()

		if err := network.EncodeXML(file, "", "  "); err != nil {
			log.Fatalf("unable to marshal output file %q: %v", settings.output, err)
		}
	default:
		if err := network.EncodeXML(os.Stdout, "", "  "); err != nil {
			log.Fatalf("unable to marshal output: %v", err)
		}
	}
}
