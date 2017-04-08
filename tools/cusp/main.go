package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var aliases = map[string]string{
	"LPOC": "3A",
}

var gps = map[string]string{
	"10417": "feruno", //# 3B-050008 @ AMBC
	"331":   "feruno", //# 3B-050057 @ CCCC
	"306":   "feruno", //# 3B-050048 @ CHHC
	"10402": "garmin", //# 3B-050019 @ CSTC
	"10409": "garmin", //# 3B-050038 @ CULC
	"10462": "garmin", //# 3A-040027 @ D03C
	"10461": "garmin", //# 3A-040007 @ D05C
	"10459": "garmin", //# 3A-040001 @ D06C
	"10444": "garmin", //# 3B-050010 @ D07C
	"10458": "garmin", //# 3B-050016 @ D08C
	"10453": "garmin", //# 3A-040023 @ D09C
	"10456": "garmin", //# 3A-040002 @ D10C
	"10405": "garmin", //# 3B-050013 @ DB1C
	"10449": "garmin", //# 3A-040009 @ DB2C
	"316":   "feruno", //# 3B-050006 @ DORC
	"227":   "feruno", //# 3A-040015 @ DSLC
	"311":   "feruno", //# 3B-050007 @ GDLC
	"206":   "feruno", //# 3A-040006 @ HORC
	"317":   "garmin", //# 3B-050050 @ HPSC
	"326":   "garmin", //# 3B-050035 @ HVSC
	"10415": "garmin", //# 3B-050030 @ KOWC
	"216":   "feruno", //# 3A-040008 @ KPOC
	"336":   "feruno", //# 3B-050058 @ LINC
	"318":   "feruno", //# 3B-050015 @ LPCC
	"226":   "feruno", //# 3A-040018 @ LPOC
	"217":   "feruno", //# 3A-040013 @ LRSC
	"10410": "garmin", //# 3b-050040 @ LSRC
	"210":   "feruno", //# 3A-040005 @ MAYC
	"330":   "feruno", //# 3B-050047 @ MSMC
	"323":   "garmin", //# 3B-050021 @ NBLC
	"313":   "feruno", //# 3B-050009 @ PEEC
	"329":   "feruno", //# 3B-050056 @ PRPC
	"307":   "garmin", //# 3B-050020 @ RHSC
	"220":   "feruno", //# 3A-040016 @ RKAC
	"211":   "feruno", //# 3A-040003 @ ROLC
	"215":   "feruno", //# 3A-040010 @ SBRC
	"10414": "garmin", //# 3B-050027 @ SCAC
	"10400": "garmin", //# 3B-050024 @ SHFC
	"325":   "garmin", //# 3B-050026 @ SHLC
	"10467": "garmin", //# 3B-050017 @ SLRC
	"327":   "garmin", //# 3B-050049 @ SMTC
	"10403": "hybrid", //# 3B-050014 @ SWNC
	"221":   "feruno", //# 3A-040024 @ TPLC
	"10411": "feruno", //# 3B-050052 @ WAKC
	"10407": "garmin", //# 3B-050031 @ WIGC
	"225":   "feruno", //# 3A-040025 @ WSFC
	"10416": "feruno", //# 3a-040022 @ WTMC

	"357": "feruno", //# MSMC - guess
	"213": "feruno", //# MSMC

	"218": "feruno", //# SHFC
	"309": "feruno", //# WTMC - guess
	"312": "feruno", //# CULC
	"315": "feruno", //# CULC
	"222": "feruno", //# MSMC
	"200": "feruno", //# AMBC
	"303": "hybrid", //# SWNC - odd as it has feruno hardware and garmin firmware

	"205": "feruno", //# CSTC
	"208": "feruno", //# LSRC
	"345": "feruno", //# SCAC
	"212": "feruno", //# WIGC
	"214": "feruno", //# LINC
}

func main() {

	var base string
	flag.StringVar(&base, "base", "../..", "delta base files")

	var output string
	flag.StringVar(&output, "output", "", "output cusp xml file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a cusp XML file from delta meta & response information\n")
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

	sites, err := buildSites(base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem loading stations %s: %v\n", base, err)
		os.Exit(1)
	}

	res, err := encodeSites(sites)
	if err != nil {
		log.Fatalf("error: unable to marshal xml: %v", err)
	}

	// output as needed ...
	switch {
	case output != "":
		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			log.Fatalf("error: unable to create directory %s: %v", filepath.Dir(output), err)
		}
		if err := ioutil.WriteFile(output, res, 0644); err != nil {
			log.Fatalf("error: unable to write file %s: %v", output, err)
		}
	default:
		os.Stdout.Write(res)
	}

}
