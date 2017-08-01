package main

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GeoNet/delta/metadb"
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

				switch installation.Datalogger.Model {
				case "Q330/3", "Q330/6", "Q330HR/6", "Q330S/3", "Q330S/6", "Q330HRS/6":
				case "BASALT", "OBSIDIAN", "CUSP3D", "CUSP3C", "Obsidian 4X Datalogger":
				default:
					continue
				}

				switch installation.Sensor.Model {
				case "FBA-ES-T", "FBA-ES-T-ISO", "FBA-ES-T-BASALT", "FBA-ES-T-OBSIDIAN", "CUSP3C", "CUSP3D":
				case "CMG-3ESP", "CMG-3ESPC", "CMG-3TB", "CMG-3TB-GN", "STS-2", "Nanometrics Trillium 120QA", "Nanometrics Trillium Compact PH TC120-PH2":
				case "L4C-3D", "L4C", "LE-3Dlite", "LE-3DliteMkII":
				default:
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

				lookup := response.Channels(func() bool {
					if installation.Sensor.Azimuth != 0.0 {
						return true
					}
					return stream.Axial
				}())
				for pin, _ := range response.Components {
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
