package resp

// Provide a stream list for a given datalogger and sensor pair
func Streams(datalogger, sensor string) []Stream {
	var streams []Stream

	// make sure we know about the sensor model - for the components
	model, ok := SensorModels[sensor]
	if !ok {
		return nil
	}

	for _, response := range Responses {
		for _, lo := range response.Dataloggers {
			for _, dataloggerModel := range lo.DataloggerList {
				if datalogger != dataloggerModel {
					continue
				}
				for _, se := range response.Sensors {
					for _, sensorModel := range se.SensorList {
						if sensor != sensorModel {
							continue
						}
						streams = append(streams, Stream{
							Datalogger: lo,
							Sensor:     se,
							Components: model.Components,
						})
					}
				}
			}
		}
	}

	return streams
}
