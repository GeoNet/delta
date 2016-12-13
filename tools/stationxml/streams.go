package main

import (
	"github.com/GeoNet/delta/resp"
)

func GetResponseStreams(datalogger, sensor string) []resp.Stream {
	var streams []resp.Stream

	for _, r := range resp.Responses {
		for _, l := range r.Dataloggers {
			for _, d := range l.Dataloggers {
				if datalogger != d {
					continue
				}
				for _, x := range r.Sensors {
					for _, b := range x.Sensors {
						if sensor == b {
							streams = append(streams, resp.Stream{
								Datalogger: l,
								Sensor:     x,
							})
						}
					}
				}
			}
		}
	}

	return streams
}
