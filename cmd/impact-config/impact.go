package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"
)

var Q = map[int]float64{
	200: 0.98829,
	100: 0.97671,
	50:  0.95395,
}

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

func (s Settings) ImpactStreams(set *meta.Set, response *resp.Resp) (map[string]Stream, error) {

	streams := make(map[string]Stream)

	// run through each station ....
	for _, station := range set.Stations() {
		if s.skip.MatchString(station.Network) {
			continue
		}

		network, ok := set.Network(station.Network)
		if !ok {
			continue
		}

		for _, site := range set.Sites() {
			if site.Station != station.Code {
				continue
			}

			for _, collection := range set.Collections(site) {
				if !s.channels.MatchString(collection.Code()) {
					continue
				}
				if collection.Component.Response == "" {
					continue
				}
				if collection.Channel.Response == "" {
					continue
				}

				if time.Now().After(collection.Span.End) {
					continue
				}

				q, ok := Q[int(collection.Stream.SamplingRate)]
				if !ok {
					continue
				}

				info, err := response.Info(collection.Component.Response, collection.Channel.Response)
				if err != nil {
					return nil, fmt.Errorf("unable to find response information for %s <=> %s (%v)",
						collection.Component.Response, collection.Channel.Response, err,
					)
				}

				key := strings.Join([]string{func() string {
					if network.External != "" {
						return network.External
					}
					return network.Code
				}(), station.Code, site.Location, collection.Code()}, "_")

				streams[key] = Stream{
					Latitude:  station.Latitude,
					Longitude: station.Longitude,
					Q:         q,
					Name:      station.Name,
					Rate:      collection.Stream.SamplingRate,
					Gain:      Float64(info.Sensitivity),
				}
			}
		}
	}

	return streams, nil
}
