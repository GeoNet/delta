package resp

// Provide a stream list for a given datalogger and sensor pair
func Streams(datalogger, sensor string) []Stream {
	var streams []Stream

	for _, r := range Responses {
		for _, l := range r.Dataloggers {
			for _, d := range l.DataloggerList {
				if datalogger != d {
					continue
				}
				for _, x := range r.Sensors {
					for _, b := range x.SensorList {
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
