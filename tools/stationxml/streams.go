package main

func GetResponseStreams(datalogger, sensor string) []Stream {
	var streams []Stream

	for _, r := range Responses {
		for _, l := range r.Dataloggers {
			for _, d := range l.Dataloggers {
				if datalogger != d {
					continue
				}
				for _, x := range r.Sensors {
					for _, b := range x.Sensors {
						if sensor == b {
							streams = append(streams, Stream{
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
