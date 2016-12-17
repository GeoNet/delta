package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"
)

// default 1st order filter factors
var Q = map[int]float64{
	200: 0.98829,
	100: 0.97671,
	50:  0.95395,
}

// the configuration results
type Stream struct {
	Source     string
	Datalogger string
	Sensor     string
	Network    string
	Station    string
	Location   string
	Channel    string
	SampleRate float64
	Gain       float64
	Q          float64
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Provide an example response configuration\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}

	var base string
	flag.StringVar(&base, "base", "../..", "base of delta files on disk")

	var expression string
	flag.StringVar(&expression, "match", "^[A-Z0-9_]+_[BH]N[A-Z0-9]$", "stream sources to match")

	flag.Parse()

	match, err := regexp.Compile(expression)
	if err != nil {
		log.Fatalf("error: unable to compile stream matching expression: %v", err)
	}

	// load the networks list
	var networkList meta.NetworkList
	if err := meta.LoadList(filepath.Join(filepath.Join(base, "network", "networks.csv")), &networkList); err != nil {
		log.Fatalf("error: unable to load networks: %v", err)
	}
	// build a map of external network codes
	networkMap := make(map[string]string)
	for _, network := range networkList {
		networkMap[network.Code] = network.External
	}

	// load the stations list
	var stationList meta.StationList
	if err := meta.LoadList(filepath.Join(filepath.Join(base, "network", "stations.csv")), &stationList); err != nil {
		log.Fatalf("error: unable to load stations: %v", err)
	}

	// build a map of station code to network codes
	stationMap := make(map[string]string)
	for _, station := range stationList {
		if network, ok := networkMap[station.Network]; ok {
			stationMap[station.Code] = network
		}
	}

	// load the streams list
	var streamList meta.StreamList
	if err := meta.LoadList(filepath.Join(filepath.Join(base, "install", "streams.csv")), &streamList); err != nil {
		log.Fatalf("error: unable to load streams: %v", err)
	}

	// load the connections list
	var connectionList meta.ConnectionList
	if err := meta.LoadList(filepath.Join(filepath.Join(base, "install", "connections.csv")), &connectionList); err != nil {
		log.Fatalf("error: unable to load connections: %v", err)
	}

	// load the installed recorder list
	var installedRecorderList meta.InstalledRecorderList
	if err := meta.LoadList(filepath.Join(filepath.Join(base, "install", "recorders.csv")), &installedRecorderList); err != nil {
		log.Fatalf("error: unable to load recorders: %v", err)
	}

	// load the installed sensors list
	var installedSensorList meta.InstalledSensorList
	if err := meta.LoadList(filepath.Join(filepath.Join(base, "install", "sensors.csv")), &installedSensorList); err != nil {
		log.Fatalf("error: unable to load recorders: %v", err)
	}

	// load the deployed dataloggers list
	var deployedDataloggerList meta.DeployedDataloggerList
	if err := meta.LoadList(filepath.Join(filepath.Join(base, "install", "dataloggers.csv")), &deployedDataloggerList); err != nil {
		log.Fatalf("error: unable to load recorders: %v", err)
	}

	streams := make(map[string]Stream)

	// run through the configured streams
	for _, stream := range streamList {
		// must be current
		if time.Now().After(stream.End) {
			continue
		}

		// do have we a filter parameter for this sampling rate
		if _, ok := Q[int(stream.SamplingRate)]; !ok {
			continue
		}

		// make sure we have a network code
		network, ok := stationMap[stream.Station]
		if !ok {
			continue
		}

		// run through the recorder lists
		for _, recorder := range installedRecorderList {

			// must be currently installed
			if time.Now().After(recorder.End) {
				continue
			}

			// at the same stream station and location
			if stream.Station != recorder.Station {
				continue
			}
			if stream.Location != recorder.Location {
				continue
			}

			// for this model pair recover the response stages
			for _, response := range resp.Streams(recorder.DataloggerModel, recorder.InstalledSensor.Model) {
				if stream.SamplingRate != response.SampleRate {
					continue
				}

				for _, channel := range response.Channels(stream.Axial) {
					source := strings.Join([]string{network, stream.Station, stream.Location, channel}, "_")
					if !match.MatchString(source) {
						continue
					}

					// add the stream
					if q, ok := Q[int(response.SampleRate)]; ok {
						streams[source] = Stream{
							Sensor:     recorder.InstalledSensor.Model,
							Datalogger: recorder.DataloggerModel,
							Source:     source,
							Network:    network,
							Station:    stream.Station,
							Location:   stream.Location,
							Channel:    channel,
							SampleRate: response.SampleRate,
							Gain:       response.Gain(),
							Q:          q,
						}
					}
				}

			}

		}

		// run through installed sensors
		for _, sensor := range installedSensorList {
			// must still be installed
			if time.Now().After(sensor.End) {
				continue
			}

			// must be at the correct station and location
			if stream.Station != sensor.Station {
				continue
			}
			if stream.Location != sensor.Location {
				continue
			}

			// run through the datalogger to sensor connections
			for _, connection := range connectionList {
				// must still be current
				if time.Now().After(connection.End) {
					continue
				}

				// must be deloyed to the right station and location
				if connection.Station != sensor.Station {
					continue
				}
				if connection.Location != sensor.Location {
					continue
				}

				// run through the deployed dataloggers
				for _, datalogger := range deployedDataloggerList {
					// must still be deployed
					if time.Now().After(datalogger.End) {
						continue
					}

					// must also be at the connection place and role
					if datalogger.Place != connection.Place {
						continue
					}
					if datalogger.Role != connection.Role {
						continue
					}

					// run through response stages for this equipment pair
					for _, response := range resp.Streams(datalogger.Model, sensor.Model) {
						if stream.SamplingRate != response.SampleRate {
							continue
						}
						for _, channel := range response.Channels(stream.Axial) {

							source := strings.Join([]string{network, stream.Station, stream.Location, channel}, "_")
							if !match.MatchString(source) {
								continue
							}

							// add current stream
							if q, ok := Q[int(response.SampleRate)]; ok {
								streams[source] = Stream{
									Sensor:     sensor.Model,
									Datalogger: datalogger.Model,
									Source:     source,
									Network:    network,
									Station:    stream.Station,
									Location:   stream.Location,
									Channel:    channel,
									SampleRate: response.SampleRate,
									Gain:       response.Gain(),
									Q:          q,
								}
							}
						}

					}
				}
			}
		}
	}

	res, err := json.MarshalIndent(streams, "", "  ")
	if err != nil {
		log.Fatalf("error: unable to produce json: %v", err)

	}

	fmt.Println(string(res))
}
