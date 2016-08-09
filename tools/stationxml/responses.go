package main

func ResponseMap() map[string]map[string][]Stream {
	resmap := make(map[string]map[string][]Stream)
	for _, r := range Responses {
		for _, l := range r.Dataloggers {
			for _, a := range l.Dataloggers {
				if _, ok := resmap[a]; !ok {
					resmap[a] = make(map[string][]Stream)
				}
				for _, x := range r.Sensors {
					for _, b := range x.Sensors {
						resmap[a][b] = append(resmap[a][b], Stream{
							Datalogger: l,
							Sensor:     x,
						})
					}
				}
			}
		}
	}

	return resmap
}
