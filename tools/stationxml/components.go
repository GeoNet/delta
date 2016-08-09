package main

func Components() map[string]map[int]SensorComponent {
	components := make(map[string]map[int]SensorComponent)

	for k, v := range SensorModels {
		if _, ok := components[k]; !ok {
			components[k] = make(map[int]SensorComponent)
		}
		for n, p := range v.Components {
			components[k][n] = p
		}
	}

	return components
}
