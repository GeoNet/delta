package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/kit/sit_delta_pb"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	//Meta Loading Code grabbed from rinexml
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var output string
	flag.StringVar(&output, "output", "output", "output directory")

	var network string
	flag.StringVar(&network, "network", "../../network", "base network directory")

	var install string
	flag.StringVar(&install, "install", "../../install", "base install directory")

	var asset string
	flag.StringVar(&asset, "asset", "../../assets", "base asset directory")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "\n")
		fmt.Fprint(os.Stderr, "Build a set of ProtoBuf files for Gloria from delta meta information\n")
		fmt.Fprint(os.Stderr, "\n")
		fmt.Fprint(os.Stderr, "Usage:\n")
		fmt.Fprint(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprint(os.Stderr, "\n")
		fmt.Fprint(os.Stderr, "Options:\n")
		fmt.Fprint(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n")
	}

	flag.Parse()

	//data from 'network' files
	//List of marks from marks.csv - gps
	var markList meta.MarkList
	if err := meta.LoadList(filepath.Join(network, "marks.csv"), &markList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load mark list: %v\n", err)
		os.Exit(-1)
	}

	//List of monuments from monuments.csv
	var monumentList meta.MonumentList
	if err := meta.LoadList(filepath.Join(network, "monuments.csv"), &monumentList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load monuments list: %v\n", err)
		os.Exit(-1)
	}
	monuments := make(map[string][]meta.Monument)
	for _, m := range monumentList {
		monuments[m.Mark] = append(monuments[m.Mark], m)
	}

	//List of mounts from mounts.csv - cameras
	var mountList meta.MountList
	if err := meta.LoadList(filepath.Join(network, "mounts.csv"), &mountList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load mounts list: %v\n", err)
		os.Exit(-1)
	}

	//List of stations from stations.csv - seismic AND tsunami
	var stationList meta.StationList
	if err := meta.LoadList(filepath.Join(network, "stations.csv"), &stationList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load stations list: %v\n", err)
		os.Exit(-1)
	}

	//List of sites from sites.csv - 'location' for stations
	var siteList meta.SiteList
	locations := make(map[string][]*sit_delta_pb.Location)
	if err := meta.LoadList(filepath.Join(network, "sites.csv"), &siteList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load sites list: %v\n", err)
		os.Exit(-1)
	}
	for _, m := range siteList {
		locations[m.Station] = append(locations[m.Station], &sit_delta_pb.Location{
			Point: &sit_delta_pb.Point{
				Datum:     m.Datum,
				Elevation: m.Elevation,
				Latitude:  m.Latitude,
				Longitude: m.Longitude,
			},
			Location: m.Location,
			Span: &sit_delta_pb.Span{
				Start: m.Start.Unix(),
				End:   m.End.Unix(),
			},
		})
	}

	//asset files -- All have the same format so just pull them all into a single big map
	assetFiles := []string{
		"antennas.csv",
		"cameras.csv",
		"dataloggers.csv",
		"metsensors.csv",
		"radomes.csv",
		"receivers.csv",
		"recorders.csv",
		"sensors.csv",
	}
	assets := make(map[string]meta.Asset)
	for _, f := range assetFiles {
		var assetList meta.AssetList
		if err := meta.LoadList(filepath.Join(asset, f), &assetList); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to load assets '%s': %v\n", f, err)
			os.Exit(-1)
		}
		for _, a := range assetList {
			assets[a.Make+a.Model+a.Serial] = a
		}
	}

	//install files -- pull in each files and turn into a map of it's mark/site to a list of matching entries
	//antennas.csv
	equipment := make(map[string][]*sit_delta_pb.Equipment_Install)
	var installedAntennaList meta.InstalledAntennaList
	if err := meta.LoadList(filepath.Join(install, "antennas.csv"), &installedAntennaList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load antenna installs: %v\n", err)
		os.Exit(-1)
	}
	for _, i := range installedAntennaList {
		equipment[i.Mark] = append(equipment[i.Mark], &sit_delta_pb.Equipment_Install{
			Equipment: &sit_delta_pb.Equipment{
				Type:         "Antenna",
				Model:        i.Model,
				AssetNumber:  assets[i.Make+i.Model+i.Serial].Number,
				SerialNumber: i.Serial,
				Manufacturer: i.Make,
				Height:       -i.Vertical,
			},
			Installed: &sit_delta_pb.Span{
				Start: i.Start.Unix(),
				End:   i.End.Unix(),
			},
		})
	}

	//cameras.csv
	var installedCameraList meta.InstalledCameraList
	if err := meta.LoadList(filepath.Join(install, "cameras.csv"), &installedCameraList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load camera installs: %v\n", err)
		os.Exit(-1)
	}
	for _, i := range installedCameraList {
		equipment[i.Mount] = append(equipment[i.Mount], &sit_delta_pb.Equipment_Install{
			Equipment: &sit_delta_pb.Equipment{
				Type:         "Camera",
				Model:        i.Model,
				AssetNumber:  assets[i.Make+i.Model+i.Serial].Number,
				SerialNumber: i.Serial,
				Manufacturer: i.Make,
				Height:       -i.Vertical,
			},
			Installed: &sit_delta_pb.Span{
				Start: i.Start.Unix(),
				End:   i.End.Unix(),
			},
		})
	}

	//connections.csv - needed to link a site to a datalogger
	var connectionList meta.ConnectionList
	if err := meta.LoadList(filepath.Join(install, "connections.csv"), &connectionList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load datalogger installs: %v\n", err)
		os.Exit(-1)
	}
	connections := make(map[string][]string)
	for _, c := range connectionList {
		//TODO - Do we need to display 'connection' info?

		arr := connections[strings.TrimSpace(c.Place+c.Role)]
		found := false
		for _, s := range arr {
			if s == c.Station {
				found = true
				break
			}
		}
		if !found {
			connections[strings.TrimSpace(c.Place+c.Role)] = append(connections[strings.TrimSpace(c.Place+c.Role)], c.Station)
		}
	}

	//dataloggers.csv
	var deployedDataloggerList meta.DeployedDataloggerList
	if err := meta.LoadList(filepath.Join(install, "dataloggers.csv"), &deployedDataloggerList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load datalogger installs: %v\n", err)
		os.Exit(-1)
	}
	for _, i := range deployedDataloggerList {
		stations := connections[strings.TrimSpace(i.Place+i.Role)]
		for _, s := range stations {
			equipment[s] = append(equipment[s], &sit_delta_pb.Equipment_Install{
				Equipment: &sit_delta_pb.Equipment{
					Type:         "Datalogger",
					Model:        i.Model,
					AssetNumber:  assets[i.Make+i.Model+i.Serial].Number,
					SerialNumber: i.Serial,
					Manufacturer: i.Make,
				},
				Installed: &sit_delta_pb.Span{
					Start: i.Start.Unix(),
					End:   i.End.Unix(),
				},
			})
		}
	}

	//metsensors.csv
	var installedMetsensorList meta.InstalledMetSensorList
	if err := meta.LoadList(filepath.Join(install, "metsensors.csv"), &installedMetsensorList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load metsensor installs: %v\n", err)
		os.Exit(-1)
	}
	for _, i := range installedMetsensorList {
		equipment[i.Mark] = append(equipment[i.Mark], &sit_delta_pb.Equipment_Install{
			Equipment: &sit_delta_pb.Equipment{
				Type:         "Metsensor",
				Model:        i.Model,
				AssetNumber:  assets[i.Make+i.Model+i.Serial].Number,
				SerialNumber: i.Serial,
				Manufacturer: i.Make,
			},
			Installed: &sit_delta_pb.Span{
				Start: i.Start.Unix(),
				End:   i.End.Unix(),
			},
		})
	}

	//recorders.csv
	var installedRecorderList meta.InstalledRecorderList
	if err := meta.LoadList(filepath.Join(install, "recorders.csv"), &installedRecorderList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load recorder installs: %v\n", err)
		os.Exit(-1)
	}
	for _, i := range installedRecorderList {
		equipment[i.Station] = append(equipment[i.Station], &sit_delta_pb.Equipment_Install{
			Equipment: &sit_delta_pb.Equipment{
				Type:         "Recorder",
				Model:        i.Model,
				AssetNumber:  assets[i.Make+i.Model+i.Serial].Number,
				SerialNumber: i.Serial,
				Manufacturer: i.Make,
				Height:       -i.Vertical,
			},
			Installed: &sit_delta_pb.Span{
				Start: i.Start.Unix(),
				End:   i.End.Unix(),
			},
		})
		for _, l := range locations[i.Station] {
			if l.Location == i.Location {
				l.GroundRelationship = -i.Vertical
			}
		}
	}

	//receivers.csv
	var deployedReceiverList meta.DeployedReceiverList
	if err := meta.LoadList(filepath.Join(install, "receivers.csv"), &deployedReceiverList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load receiver installs: %v\n", err)
		os.Exit(-1)
	}
	for _, i := range deployedReceiverList {
		equipment[i.Mark] = append(equipment[i.Mark], &sit_delta_pb.Equipment_Install{
			Equipment: &sit_delta_pb.Equipment{
				Type:         "Receiver",
				Model:        i.Model,
				AssetNumber:  assets[i.Make+i.Model+i.Serial].Number,
				SerialNumber: i.Serial,
				Manufacturer: i.Make,
			},
			Installed: &sit_delta_pb.Span{
				Start: i.Start.Unix(),
				End:   i.End.Unix(),
			},
		})
	}

	//radomes.csv
	var installedRadomeList meta.InstalledRadomeList
	if err := meta.LoadList(filepath.Join(install, "radomes.csv"), &installedRadomeList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load radome installs: %v\n", err)
		os.Exit(-1)
	}
	for _, i := range installedRadomeList {
		equipment[i.Mark] = append(equipment[i.Mark], &sit_delta_pb.Equipment_Install{
			Equipment: &sit_delta_pb.Equipment{
				Type:         "Radome",
				Model:        i.Model,
				AssetNumber:  assets[i.Make+i.Model+i.Serial].Number,
				SerialNumber: i.Serial,
				Manufacturer: i.Make,
			},
			Installed: &sit_delta_pb.Span{
				Start: i.Start.Unix(),
				End:   i.End.Unix(),
			},
		})
	}

	//sensors.csv
	var installedSensorList meta.InstalledSensorList
	if err := meta.LoadList(filepath.Join(install, "sensors.csv"), &installedSensorList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load sensor installs: %v\n", err)
		os.Exit(-1)
	}
	for _, i := range installedSensorList {
		equipment[i.Station] = append(equipment[i.Station], &sit_delta_pb.Equipment_Install{
			Equipment: &sit_delta_pb.Equipment{
				Type:         "Sensor",
				Model:        i.Model,
				AssetNumber:  assets[i.Make+i.Model+i.Serial].Number,
				SerialNumber: i.Serial,
				Manufacturer: i.Make,
				Height:       -i.Vertical,
			},
			Installed: &sit_delta_pb.Span{
				Start: i.Start.Unix(),
				End:   i.End.Unix(),
			},
		})
		for _, l := range locations[i.Station] {
			if l.Location == i.Location {
				l.GroundRelationship = -i.Vertical
			}
		}
	}

	for _, m := range markList {

		im := make([]*sit_delta_pb.InstalledMonument, 0)

		list := monuments[m.Code]
		var currentMon *meta.Monument
		for _, l := range list {
			newMonument := sit_delta_pb.InstalledMonument{
				Span: &sit_delta_pb.Span{
					Start: l.Start.Unix(),
					End:   l.End.Unix(),
				},
				Monument: &sit_delta_pb.Monument{
					DomesNumber: l.DomesNumber,
					Height:      l.GroundRelationship,
				},
			}
			im = append(im, &newMonument)
			if l.End.Unix() > time.Now().Unix() {
				currentMon = &l
			}
		}

		mark := sit_delta_pb.Mark{
			InstalledMonument: im,
			Point: &sit_delta_pb.Point{
				Longitude: m.Longitude,
				Latitude:  m.Latitude,
				Elevation: m.Elevation,
				Datum:     m.Datum,
			},
		}

		site_pb := sit_delta_pb.Site{
			Code:               m.Code,
			Span:               &sit_delta_pb.Span{Start: m.Start.Unix(), End: m.End.Unix()},
			Network:            m.Network,
			Mark:               &mark,
			Point:              &sit_delta_pb.Point{Longitude: m.Longitude, Latitude: m.Latitude, Elevation: m.Elevation, Datum: m.Datum},
			GroundRelationship: 0,
			EquipmentInstalls:  equipment[m.Code],
		}
		if currentMon != nil {
			site_pb.GroundRelationship = currentMon.GroundRelationship
		}

		b, err := proto.Marshal(&site_pb)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to marsh protobuf: %v\n", err)
			os.Exit(-1)
		}

		pbfile := filepath.Join(output, strings.ToUpper(m.Code)+".pb")
		if err := os.MkdirAll(filepath.Dir(pbfile), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to create dir: %v\n", err)
			os.Exit(-1)
		}
		if err := ioutil.WriteFile(pbfile, b, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to write file: %v\n", err)
			os.Exit(-1)
		}
		if verbose {
			out_json, _ := json.MarshalIndent(site_pb, "", "  ")
			ioutil.WriteFile(filepath.Join(output, strings.ToUpper(m.Code)+".json"), []byte(out_json), 0644)
		}
	}

	for _, m := range stationList {
		site_pb := sit_delta_pb.Site{
			Code:               m.Code,
			Span:               &sit_delta_pb.Span{Start: m.Start.Unix(), End: m.End.Unix()},
			Network:            m.Network,
			Point:              &sit_delta_pb.Point{Longitude: m.Longitude, Latitude: m.Latitude, Elevation: m.Elevation, Datum: m.Datum},
			GroundRelationship: 0,
			EquipmentInstalls:  equipment[m.Code],
			Locations:          locations[m.Code],
		}
		b, err := proto.Marshal(&site_pb)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to marsh protobuf: %v\n", err)
			os.Exit(-1)
		}

		pbfile := filepath.Join(output, strings.ToUpper(m.Code)+".pb")
		if err := os.MkdirAll(filepath.Dir(pbfile), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to create dir: %v\n", err)
			os.Exit(-1)
		}
		if err := ioutil.WriteFile(pbfile, b, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to write file: %v\n", err)
			os.Exit(-1)
		}
		if verbose {
			out_json, _ := json.MarshalIndent(site_pb, "", "  ")
			ioutil.WriteFile(filepath.Join(output, strings.ToUpper(m.Code)+".json"), []byte(out_json), 0644)
		}
	}

	for _, m := range mountList {
		site_pb := sit_delta_pb.Site{
			Code:              m.Code,
			Span:              &sit_delta_pb.Span{Start: m.Start.Unix(), End: m.End.Unix()},
			Network:           m.Network,
			Point:             &sit_delta_pb.Point{Longitude: m.Longitude, Latitude: m.Latitude, Elevation: m.Elevation, Datum: m.Datum},
			EquipmentInstalls: equipment[m.Code],
		}

		b, err := proto.Marshal(&site_pb)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to marsh protobuf: %v\n", err)
			os.Exit(-1)
		}

		pbfile := filepath.Join(output, strings.ToUpper(m.Code)+".pb")
		if err := os.MkdirAll(filepath.Dir(pbfile), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to create dir: %v\n", err)
			os.Exit(-1)
		}

		if err := ioutil.WriteFile(pbfile, b, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to write file: %v\n", err)
			os.Exit(-1)
		}
		if verbose {
			out_json, _ := json.MarshalIndent(site_pb, "", "  ")
			ioutil.WriteFile(filepath.Join(output, strings.ToUpper(m.Code)+".json"), []byte(out_json), 0644)
		}
	}
}
