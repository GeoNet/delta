package delta_test

import (
	"testing"
	"github.com/GeoNet/delta/meta"
	"log"
	"time"
	"github.com/GeoNet/delta/resp"
	"strings"
	"regexp"
	"github.com/GeoNet/delta"
	"os"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
)

// default 1st order filter factors
var Q = map[int]float64{
	200: 0.98829,
	100: 0.97671,
	50:  0.95395,
}

// TestShakeNZCfg generates protobuf config for shakenz-slink
func TestShakeNZCfg(t *testing.T) {
	match, err := regexp.Compile("^[A-Z0-9_]+_[BH]N[A-Z0-9]$")
	if err != nil {
		t.Fatal(err)
	}

	var networkList meta.NetworkList
	if err := meta.LoadList("../network/networks.csv", &networkList); err != nil {
		log.Fatalf("error: unable to load networks: %v", err)
	}

	networkMap := make(map[string]string)
	for _, network := range networkList {
		networkMap[network.Code] = network.External
	}

	var stationList meta.StationList
	if err := meta.LoadList("../network/stations.csv", &stationList); err != nil {
		log.Fatalf("error: unable to load stations: %v", err)
	}

	// build a map of station code to network codes
	stationMap := make(map[string]string)
	for _, station := range stationList {
		if network, ok := networkMap[station.Network]; ok {
			stationMap[station.Code] = network
		}
	}

	var streamList meta.StreamList
	if err := meta.LoadList("../install/streams.csv", &streamList); err != nil {
		log.Fatalf("error: unable to load streams: %v", err)
	}

	var connectionList meta.ConnectionList
	if err := meta.LoadList("../install/connections.csv", &connectionList); err != nil {
		log.Fatalf("error: unable to load connections: %v", err)
	}

	var installedRecorderList meta.InstalledRecorderList
	if err := meta.LoadList("../install/recorders.csv", &installedRecorderList); err != nil {
		log.Fatalf("error: unable to load recorders: %v", err)
	}

	var installedSensorList meta.InstalledSensorList
	if err := meta.LoadList("../install/sensors.csv", &installedSensorList); err != nil {
		log.Fatalf("error: unable to load recorders: %v", err)
	}

	var deployedDataloggerList meta.DeployedDataloggerList
	if err := meta.LoadList("../install/dataloggers.csv", &deployedDataloggerList); err != nil {
		log.Fatalf("error: unable to load recorders: %v", err)
	}

	var streams delta.ShakeNZStreams
	streams.Streams = make(map[string]*delta.ShakeNZStream)

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
						streams.Streams[source] = &delta.ShakeNZStream{
							Sensor:     recorder.InstalledSensor.Model,
							Datalogger: recorder.DataloggerModel,
							StreamId:     source,
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
								streams.Streams[source] = &delta.ShakeNZStream{
									Sensor:     sensor.Model,
									Datalogger: datalogger.Model,
									StreamId:     source,
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

	// get a map of station and site location then we can add longitude and latitude to the config.
	s, err := stations()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range streams.Streams {
		if _, ok := s.Stations[v.Station]; !ok {
			t.Errorf("no station for stream %s can't find longitude latitude", k)
		}
		if _, ok := s.Stations[v.Station].Sites[v.Location]; !ok {
			t.Errorf("no site for stream %s can't find longitude latitude", k)
		}

		v.Longitude = s.Stations[v.Station].Sites[v.Location].Point.Longitude
		v.Latitude = s.Stations[v.Station].Sites[v.Location].Point.Latitude

		v.Source = v.Network + "." + v.Station + "." + v.Location

		if len(v.Channel) != 3 {
			t.Errorf("chan not len 3 for %s", k)
		}

		// TODO should we also use components 1,2,3?
		switch strings.ToUpper(string(v.Channel[2])) {
		case "Z":
			v.Vertical = true
			streams.Streams[k] = v
		case "N", "E":
			v.Horizontal = true
			streams.Streams[k] = v
		default:
			delete(streams.Streams, k)
		}
	}

	if len(streams.Streams) == 0 {
		t.Fatal("got zero length streams for shakenz config.")
	}

	b, err := proto.Marshal(&streams)
	if err != nil {
		t.Fatal(err)

	}

	if err := os.MkdirAll(apiDir + "/config", 0777); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(apiDir + "/config/shakenz-slink.pb", b, 0644); err != nil {
		t.Fatal(err)
	}
}