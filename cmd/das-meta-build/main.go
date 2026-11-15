package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// "path/filepath"
	"strings"
	// "text/template"
	"time"

	"github.com/GeoNet/delta"
	// "delta"

	dasroot "github.com/GeoNet/delta/internal/das-json"
	"github.com/GeoNet/delta/meta"
)

const ClockDrift = 0.0001

const (
	networkRe     = `[A-Z0-9]+`
	locationRe    = `[A-Z0-9]+`
	acquisitionRe = `[A-Z0-9]+`
)

// default response frequency values

type Settings struct {
	verbose bool // output operational info
	debug   bool // output more operational info

	base string // base directory of delta files on disk
	resp string // base directory for response xml files on disk

	version string // create a specific json version - connect to json schema download
	// create      bool   // add a root XML Created entry
	// corrections bool   // add calculated and applied response delays and corrections

	// source string // stationxml source
	// sender string // stationxml sender
	// module string // stationxml module

	// external Matcher // regexp selection of external networks
	// exclude  Matcher // regexp selection of networks to exclude
	network     Matcher // regexp selection of networks
	acquisition Matcher
	// station  Matcher // regexp selection of stations
	location Matcher // regexp selection of locations
	// channel  Matcher // regexp selection of channels
	ignore string // list of stations to skip

	// single    bool   // produce single station xml files
	directory string // where to store das-json
	// template  string // how to name the single station xml files
	// purge     bool   // remove unknown single xml files

	output  string // output xml file, use - for stdout
	changed bool   // only update existing file if a change is detected

	active       bool // only output stations with active channels
	installation bool // set station times based on installation dates
	operational  bool // only output operational channels

	offset time.Duration // provide a recently closed window for operational only requests")
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a DAS json file from delta meta\n")
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

	flag.BoolVar(&settings.verbose, "verbose", false, "add operational info")
	flag.BoolVar(&settings.debug, "debug", false, "add extra operational info")

	flag.StringVar(&settings.base, "base", "", "base directory of delta files on disk")
	flag.StringVar(&settings.resp, "resp", "", "base directory for response xml files on disk")

	flag.StringVar(&settings.version, "version", "2.0", "create a specific DAS-json version")
	// flag.BoolVar(&settings.create, "create", false, "add a root XML \"Created\" entry")
	// flag.BoolVar(&settings.corrections, "corrections", false, "add calculated and applied response delays and corrections")

	// flag.StringVar(&settings.source, "source", "GeoNet", "stationxml source")
	// flag.StringVar(&settings.sender, "sender", "WEL(GNS)", "stationxml sender")
	// flag.StringVar(&settings.module, "module", "Delta", "stationxml module")

	// flag.TextVar(&settings.external, "external", MustMatcher(externalRe), "regexp selection of external networks")
	// flag.TextVar(&settings.exclude, "exclude", MustMatcher(excludeRe), "regexp selection of networks to exclude")
	flag.TextVar(&settings.network, "networkCode", MustMatcher(networkRe), "regexp selection of network codes")
	flag.TextVar(&settings.acquisition, "acquisistion", MustMatcher(acquisitionRe), "regexp selection of acquisistions")
	// flag.TextVar(&settings.station, "station", MustMatcher(stationRe), "regexp selection of stations")
	flag.TextVar(&settings.location, "location", MustMatcher(locationRe), "regexp selection of locations")
	// flag.TextVar(&settings.channel, "channel", MustMatcher(channelRe), "regexp selection of channels")
	flag.StringVar(&settings.ignore, "ignore", "", "list of stations to skip")

	// flag.BoolVar(&settings.single, "single", false, "produce single station xml files")
	flag.StringVar(&settings.directory, "directory", "xml", "where to store station xml files")
	// flag.StringVar(&settings.template, "template", "station_{{.ExternalCode}}_{{.StationCode}}.xml", "how to name the single station xml files")
	// flag.BoolVar(&settings.purge, "purge", false, "remove unknown single xml files")
	flag.StringVar(&settings.output, "output", "", "output xml file, use \"-\" for stdout")
	flag.BoolVar(&settings.changed, "changed", false, "only update existing file if a change is detected")

	flag.BoolVar(&settings.installation, "installation", false, "set station times based on installation dates")
	flag.BoolVar(&settings.operational, "operational", false, "only output operational channels")
	flag.BoolVar(&settings.active, "active", false, "only output stations with active channels")
	flag.DurationVar(&settings.offset, "operational-offset", 0, "provide a recently closed window for operational only requests")

	// flag.Func("freq", "response frequency (e.g B:1.0)", func(s string) error {
	// 	freq, err := NewFrequency(s)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	freqs.Set(freq.Prefix, freq.Value)
	// 	return nil
	// })

	flag.Parse()

	// set recovers the delta tables
	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatal(err)
	}

	// simple skip list of stations
	skip := make(map[string]interface{})
	for _, s := range strings.Split(settings.ignore, ",") {
		if v := strings.TrimSpace(s); v != "" {
			skip[v] = true
		}
	}

	// remember individual stations in case of single file output
	// var singles []string

	installed := make(map[string]interface{})
	for _, das := range set.Das() {
		installed[das.Location] = true
	}

	// find a map of stations that match
	stns := make(map[string]meta.DASStation)
	for _, d := range set.Das() {
		if !settings.location.MatchString(d.Location) {
			if settings.debug {
				log.Printf("debug: skip das %q, doesn't match location regexp", d.Location)
			}
			continue
		}
		if !settings.network.MatchString(d.Network) {
			if settings.debug {
				log.Printf("debug: skip station %q, doesn't match network code regexp", d.Location)
			}
			continue
		}
		if _, ok := skip[d.Location]; ok {
			if settings.debug {
				log.Printf("debug: skip station %q, matches skip list", d.Location)
			}
			continue
		}
		if settings.debug {
			log.Printf("debug: add station %q from network %q", d.Location, d.Network)
		}
		stns[d.Location] = d
	}
	var installedDas []dasroot.DASoverview
	for _, stn := range stns {
		// for i := 0; i < reflect.ValueOf(stn).NumField(); i++ {
		// 	fmt.Print(reflect.TypeOf(stn).Field(i).Name, " | ", reflect.ValueOf(stn).Field(i), "\n")
		// }
		var interrogators []*dasroot.Interrogator
		// set.InstalledInterrogators()
		for _, i := range set.InstalledInterrogators() {
			if i.NetworkCode == stn.Network {

				var acquisitions []*dasroot.Acquisition
				for _, a := range collectMatching(set.DasAcquisitions(), "InterrogatorID", i.InterrogatorID, i.Span) {
					// fmt.Print("\n\nMATCHING!!!!!!!!!!!!!\n\n")

					var channel_groups []*dasroot.ChannelGroup
					// for _, c := range set.DasChannelGroups() {
					// 	fmt.Print(a.AcquisitionID, "|", c.AcquisitionID, "\n\n")
					// }
					for _, c := range collectMatching(set.DasChannelGroups(), "AcquisitionID", a.AcquisitionID, a.Span) {
						var channels dasroot.Channels
						// fmt.Print("Here\n")
						for _, channel := range collectMatching(set.DasChannels(), "ChannelGroupID", c.ChannelGroupID, c.Span) {
							// fmt.Printf("Channel found %s %s\n", i, channel.ChannelID)
							channels.AppendChannel(channel)
						}
						// fmt.Print("Channels:", channels, "\n")

						channel_groups = append(channel_groups, &dasroot.ChannelGroup{
							CableID:                  c.CableID,
							ChannelGroupID:           c.ChannelGroupID,
							Channels:                 &channels,
							CoordinateGenerationDate: c.CoordinateGenerationDate,
							CoordinateSystem:         c.CoordinateSystem,
							FiberID:                  c.FiberID,
							FirstUsableChannelID:     c.FirstUsableChannelID,
							LastUsableChannelID:      c.LastUsableChannelID,
							LocationMethod:           c.LocationMethod,
							ReferenceFrame:           c.ReferenceFrame,
							UncertaintyInDepth:       c.UncertaintyInDepth,
							UncertaintyInDip:         c.UncertaintyInDip,
							UncertaintyInElevation:   c.UncertaintyInElevation,
							UncertaintyInStrike:      c.UncertaintyInStrike,
							UncertaintyInXCoordinate: c.UncertaintyInXCoordinate,
							UncertaintyInYCoordinate: c.UncertaintyInYCoordinate,
						})
					}

					acquisitions = append(acquisitions, &dasroot.Acquisition{
						AcquisitionID:           a.AcquisitionID,
						AcquisitionSampleRate:   a.SampleRate,
						GaugeLength:             a.GaugeLength,
						NativeHeaders:           a.NativeHeaders,
						NumberOfChannels:        a.NumberOfChannels,
						PulseRate:               a.PulseRate,
						PulseWidth:              a.PulseWidth,
						SpatialSamplingInterval: a.SpatialSamplingInterval,
						UnitOfMeasure:           a.UnitOfMeasure,
						AcquisitionStartTime:    a.Start,
						AcquisitionEndTime:      a.End,
						ChannelGroups:           channel_groups,
						// ScaleFactor:                 a,
						// Comment:                     a.Comment,
					})
				}

				interrogators = append(interrogators, &dasroot.Interrogator{
					InterrogatorID: i.InterrogatorID,
					Manufacturer:   i.Manufacturer,
					Model:          i.Model,
					SerialNumber:   i.SerialNumber,
					Location:       i.NetworkCode,
					Acquisitions:   acquisitions,
				})
			}
		}
		var cables []*dasroot.Cable
		for _, c := range set.InstalledCables() {

			if c.NetworkCode == stn.Network {
				// fmt.Print("MATCH!\n")
				var fibers []*dasroot.Fiber
				// for _, f := range set.InstalledFibers() {
				// 	fmt.Print("FIBER ", f.CableID, " | ", c.CableID, "\n")
				// }

				for _, f := range collectMatching(set.InstalledFibers(), "CableID", c.CableID, c.Span) {
					// print("\nFIBER\n\n")
					fibers = append(fibers, &dasroot.Fiber{
						// Comment                    string
						FiberStartLocation:     f.StartLocation,
						FiberEndLocation:       f.EndLocation,
						FiberGeometry:          f.Geometry,
						FiberID:                f.FiberID,
						FiberMode:              f.Mode,
						FiberOneWayAttenuation: f.OneWayAttenuation,
						FiberOpticLength:       f.OpticLength,
						FiberRefractionIndex:   f.RefractionIndex,
						FiberWindingAngle:      f.WindingAngle,
					})
				}

				cables = append(cables, &dasroot.Cable{
					CableBoundingBox:             c.BoundingBox,
					CableCharacteristics:         c.Characteristics,
					CableEnvironment:             c.Environment,
					CableID:                      c.CableID,
					CableInstallationDate:        c.Start,
					CableInstallationEnvironment: c.InstallationEnvironment,
					CableModel:                   c.Model,
					CableOutsideDiameter:         c.OutsideDiameter,
					CableOwner:                   c.Owner,
					CableRemovalDate:             c.End,
					Fibers:                       fibers,
					// Comment:                      c.Comment,
				})
			}
		}
		// fmt.Print("###", *cables[0])

		installedDas = append(installedDas, dasroot.DASoverview{
			Code:                  stn.Network,
			Country:               stn.Country,
			Location:              stn.Location,
			PointOfContact:        stn.PointOfContact,
			PointOfContactAddress: stn.PointOfContactAddress,
			PointOfContactEmail:   stn.PointOfContactEmail,
			StartDate:             stn.Start,
			EndDate:               stn.End,
			PrincipalInvestigator: dasroot.PrincipalInvestigator{
				Name:    stn.PrincipalInvestigator.Name,
				Email:   stn.PrincipalInvestigator.Email,
				Address: stn.PrincipalInvestigator.Address,
			},
			Interrogators: interrogators,
			Cables:        cables,
		})

	}
	root := dasroot.Root{
		DAS: installedDas,
	}

	var schema string
	if settings.version == "2.0" {
		schema = "https://www.fdsn.org/schemas/DAS-Metadata-FDSN/2.0"
	} else {
		log.Fatalf("Version %s not supported", settings.version)
	}

	if err := root.WriteFile(settings.output, settings.version, schema); err != nil {
		log.Fatalf("unable to encode response %s: %v", settings.output, err)
	}

}
