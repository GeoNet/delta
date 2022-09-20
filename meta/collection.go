package meta

import (
	"sort"
)

// Collection describes the period where a sensor and a datalogger are co-located at a site with the associated streams.
type Collection struct {
	Span

	Stream
	Channel
	Component

	InstalledSensor
	DeployedDatalogger

	Gains []Gain

	SensorCalibrations     []Calibration
	DataloggerCalibrations []Calibration
}

// Less compares whether one Collection will sort before another.
func (c Collection) Less(collection Collection) bool {
	switch {
	case c.InstalledSensor.Station < collection.InstalledSensor.Station:
		return true
	case c.InstalledSensor.Station > collection.InstalledSensor.Station:
		return false
	case c.InstalledSensor.Location < collection.InstalledSensor.Location:
		return true
	case c.InstalledSensor.Location > collection.InstalledSensor.Location:
		return false
	case c.Component.Number < collection.Component.Number:
		return true
	case c.Component.Number > collection.Component.Number:
		return false
	case c.Channel.Number < collection.Channel.Number:
		return true
	case c.Channel.Number > collection.Channel.Number:
		return false
	case c.Span.Start.Before(collection.Span.Start):
		return true
	case c.Span.Start.After(collection.Span.Start):
		return false
	case c.Stream.SamplingRate > collection.Stream.SamplingRate:
		return true
	default:
		return false
	}
}

// Collections decodes the stored sensor and datalogger installation
// times and builds a slice of overlapping time spans for the given site.
func (s *Set) Collections(site Site) []Collection {
	var collections []Collection

	for _, recorder := range s.InstalledRecorders() {
		if recorder.Station != site.Station {
			continue
		}
		if recorder.Location != site.Location {
			continue
		}

		for _, stream := range s.Streams() {
			if stream.Station != site.Station {
				continue
			}
			if stream.Location != site.Location {
				continue
			}

			span, ok := recorder.Span.Extent(stream.Span)
			if !ok {
				continue
			}

			for _, component := range s.Components() {
				if recorder.InstalledSensor.Make != component.Make {
					continue
				}
				if recorder.InstalledSensor.Model != component.Model {
					continue
				}

				for _, channel := range s.Channels() {
					if recorder.Make != channel.Make {
						continue
					}

					if recorder.DataloggerModel != channel.Model {
						continue
					}
					if stream.SamplingRate != channel.SamplingRate {
						continue
					}

					var gains []Gain
					for _, g := range s.Gains() {
						if g.Station != stream.Station {
							continue
						}
						if g.Location != stream.Location {
							continue
						}
						if g.Subsource != component.Subsource {
							continue
						}
						if !span.Overlaps(g.Span) {
							continue
						}
						gains = append(gains, g)
					}
					sort.Slice(gains, func(i, j int) bool {
						return gains[i].Span.Start.Before(gains[j].Span.Start)
					})

					var sensors []Calibration
					for _, c := range s.Calibrations() {
						if c.Make != recorder.InstalledSensor.Make {
							continue
						}
						if c.Model != recorder.InstalledSensor.Model {
							continue
						}
						if c.Serial != recorder.InstalledSensor.Serial {
							continue
						}
						if !span.Overlaps(c.Span) {
							continue
						}
						sensors = append(sensors, c)
					}
					sort.Slice(sensors, func(i, j int) bool {
						return sensors[i].Span.Start.Before(sensors[j].Span.Start)
					})

					var dataloggers []Calibration
					for _, c := range s.Calibrations() {
						if c.Make != recorder.InstalledSensor.Make {
							continue
						}
						if c.Model != recorder.DataloggerModel {
							continue
						}
						if c.Serial != recorder.InstalledSensor.Serial {
							continue
						}
						if !span.Overlaps(c.Span) {
							continue
						}
						dataloggers = append(dataloggers, c)
					}
					sort.Slice(dataloggers, func(i, j int) bool {
						return dataloggers[i].Span.Start.Before(dataloggers[j].Span.Start)
					})

					collections = append(collections, Collection{
						InstalledSensor: recorder.InstalledSensor,
						DeployedDatalogger: DeployedDatalogger{
							Install: Install{
								Equipment: Equipment{
									Make:   recorder.InstalledSensor.Make,
									Model:  recorder.DataloggerModel,
									Serial: recorder.InstalledSensor.Serial,
								},
								Span: Span{
									Start: recorder.Start,
									End:   recorder.End,
								},
							},
						},
						Stream:                 stream,
						Gains:                  gains,
						SensorCalibrations:     sensors,
						DataloggerCalibrations: dataloggers,
						Channel:                channel,
						Component:              component,
						Span:                   span,
					})
				}
			}
		}
	}

	for _, connection := range s.Connections() {
		if connection.Station != site.Station {
			continue
		}
		if connection.Location != site.Location {
			continue
		}

		for _, sensor := range s.InstalledSensors() {
			if sensor.Station != site.Station {
				continue
			}
			if sensor.Location != site.Location {
				continue
			}

			for _, component := range s.Components() {
				if sensor.Make != component.Make {
					continue
				}
				if sensor.Model != component.Model {
					continue
				}

				for _, datalogger := range s.DeployedDataloggers() {
					if datalogger.Place != connection.Place {
						continue
					}
					if datalogger.Role != connection.Role {
						continue
					}

					span, ok := connection.Span.Extent(sensor.Span, datalogger.Span)
					if !ok {
						continue
					}

					for _, stream := range s.Streams() {
						if stream.Station != site.Station {
							continue
						}
						if stream.Location != site.Location {
							continue
						}

						span, ok := span.Extent(stream.Span)
						if !ok {
							continue
						}

						for _, channel := range s.Channels() {
							if datalogger.Make != channel.Make {
								continue
							}
							if datalogger.Model != channel.Model {
								continue
							}
							if component.Number+connection.Number < channel.Number {
								continue
							}

							if stream.SamplingRate != channel.SamplingRate {
								continue
							}

							var gains []Gain
							for _, g := range s.Gains() {
								if g.Station != stream.Station {
									continue
								}
								if g.Location != stream.Location {
									continue
								}
								if g.Subsource != component.Subsource {
									continue
								}
								if !span.Overlaps(g.Span) {
									continue
								}
								gains = append(gains, g)
							}
							sort.Slice(gains, func(i, j int) bool {
								return gains[i].Span.Start.Before(gains[j].Span.Start)
							})

							var sensors []Calibration
							for _, c := range s.Calibrations() {
								if c.Make != sensor.Make {
									continue
								}
								if c.Model != sensor.Model {
									continue
								}
								if c.Serial != sensor.Serial {
									continue
								}
								if c.Number != component.Number {
									continue
								}
								if !span.Overlaps(c.Span) {
									continue
								}
								sensors = append(sensors, c)
							}
							sort.Slice(sensors, func(i, j int) bool {
								return sensors[i].Span.Start.Before(sensors[j].Span.Start)
							})

							var dataloggers []Calibration
							for _, c := range s.Calibrations() {
								if c.Make != datalogger.Make {
									continue
								}
								if c.Model != datalogger.Model {
									continue
								}
								if c.Serial != datalogger.Serial {
									continue
								}
								if c.Number != channel.Number {
									continue
								}
								if !span.Overlaps(c.Span) {
									continue
								}
								dataloggers = append(dataloggers, c)
							}
							sort.Slice(dataloggers, func(i, j int) bool {
								return dataloggers[i].Span.Start.Before(dataloggers[j].Span.Start)
							})

							collections = append(collections, Collection{
								InstalledSensor:        sensor,
								DeployedDatalogger:     datalogger,
								Stream:                 stream,
								Gains:                  gains,
								SensorCalibrations:     sensors,
								DataloggerCalibrations: dataloggers,
								Channel:                channel,
								Component:              component,
								Span:                   span,
							})

						}
					}
				}
			}
		}
	}

	sort.Slice(collections, func(i, j int) bool {
		return collections[i].Less(collections[j])
	})

	return collections
}
