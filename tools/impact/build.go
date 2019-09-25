package main

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/resp"
)

var Q = map[int]float64{
	200: 0.98829,
	100: 0.97671,
	50:  0.95395,
}

type Installations []metadb.Installation

func (t Installations) Len() int           { return len(t) }
func (t Installations) Swap(a, b int)      { t[a], t[b] = t[b], t[a] }
func (t Installations) Less(a, b int) bool { return t[a].Start.Before(t[b].Start) }

type Float64 float64

func (f Float64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 64)), nil
}

type Stream struct {
	Longitude float64 `json:"Longitude"`
	Gain      Float64 `json:"Gain"`
	Q         float64 `json:"Q"`
	Rate      float64 `rate:"Rate"`
	Name      string  `rate:"Name"`
	Latitude  float64 `json:"Latitude"`
}

func buildStreams(base, channels string) (map[string]Stream, error) {

	// load delta meta helper
	db := metadb.NewMetaDB(base)

	// load station details
	stations, err := db.Stations()
	if err != nil {
		return nil, err
	}

	streams := make(map[string]Stream)

	match, err := regexp.Compile(channels)
	if err != nil {
		return nil, err
	}

	// run through each station ....
	for _, station := range stations {
		network, err := db.Network(station.Network)
		if err != nil {
			return nil, err
		}
		if network == nil {
			continue
		}

		installations, err := db.Installations(station.Code)
		if err != nil {
			return nil, err
		}

		sort.Sort(Installations(installations))

		for _, installation := range installations {
			if time.Now().After(installation.End) {
				continue
			}
			for _, response := range resp.Streams(installation.Datalogger.Model, installation.Sensor.Model) {
				q, ok := Q[int(response.Datalogger.SampleRate)]
				if !ok {
					continue
				}

				if !isBlessedDatalogger(installation.Datalogger.Model) {
					continue
				}

				if !isBlessedSensor(installation.Sensor.Model) {
					continue
				}

				stream, err := db.StationLocationSamplingRateStartStream(
					station.Code,
					installation.Location,
					response.Datalogger.SampleRate,
					installation.Start)
				if err != nil {
					return nil, err
				}
				if stream == nil {
					continue
				}
				if time.Now().After(stream.End) {
					continue
				}

				lookup := response.Channels(func() string {
					if installation.Sensor.Azimuth != 0.0 {
						return "true"
					}
					return stream.Axial
				}())
				for pin := range response.Components {
					if !(pin < len(lookup)) {
						continue
					}
					channel := lookup[pin]
					if !match.MatchString(channel) {
						continue
					}

					key := strings.Join([]string{func() string {
						if network.External != "" {
							return network.External
						}
						return network.Code
					}(), station.Code, installation.Location, channel}, "_")

					streams[key] = Stream{
						Latitude:  station.Latitude,
						Longitude: station.Longitude,
						Q:         q,
						Name:      station.Name,
						Rate:      response.Datalogger.SampleRate,
						Gain:      Float64(response.Gain()),
					}

				}
			}
		}
	}

	return streams, nil
}
