package meta

import (
	"sort"
)

// Correction contains adjustments to a Sensor and Datalogger instrument details to account for installation settings and changes.
type Correction struct {
	Span

	Polarity              *Polarity
	Gain                  *Gain
	Preamp                *Preamp
	Telemetry             *Telemetry
	Timing                *Timing
	SensorCalibration     *Calibration
	DataloggerCalibration *Calibration
}

// Corrections returns a slice of Correction values for a given Collection.
func (set *Set) Corrections(coll Collection) []Correction {
	var corrections []Correction

	for _, polarity := range set.PolarityCorrections(coll) {
		span, ok := coll.Span.Extent(polarity.Span)
		if !ok {
			continue
		}

		for _, gain := range set.GainCorrections(coll) {
			span, ok := span.Extent(gain.Span)
			if !ok {
				continue
			}

			for _, preamp := range set.PreampCorrections(coll) {
				span, ok := span.Extent(preamp.Span)
				if !ok {
					continue
				}

				for _, telemetry := range set.TelemetryCorrections(coll) {
					span, ok := span.Extent(telemetry.Span)
					if !ok {
						continue
					}

					for _, sensor := range set.SensorCalibrationCorrections(coll) {
						span, ok := span.Extent(sensor.Span)
						if !ok {
							continue
						}

						for _, datalogger := range set.DataloggerCalibrationCorrections(coll) {
							span, ok := span.Extent(datalogger.Span)
							if !ok {
								continue
							}

							for _, timing := range set.TimingCorrections(coll) {
								span, ok := span.Extent(timing.Span)
								if !ok {
									continue
								}

								corrections = append(corrections, Correction{
									Span:                  span,
									Polarity:              polarity.Polarity,
									Gain:                  gain.Gain,
									Preamp:                preamp.Preamp,
									Telemetry:             telemetry.Telemetry,
									Timing:                timing.Timing,
									SensorCalibration:     sensor.SensorCalibration,
									DataloggerCalibration: datalogger.DataloggerCalibration,
								})
							}
						}
					}
				}
			}
		}
	}

	return corrections
}

// TimingCorrections returns a slice of Correction values to account for any changes in Timing.
func (s *Set) TimingCorrections(coll Collection) []Correction {
	var timings []Timing

	// build a slice of associated timings
	for _, t := range s.Timings() {
		if t.Station != coll.Stream.Station {
			continue
		}
		if t.Location != coll.Stream.Location {
			continue
		}
		if !t.Span.Overlaps(coll.Span) {
			continue
		}

		timings = append(timings, t)
	}

	sort.Slice(timings, func(i, j int) bool {
		return timings[i].Span.Start.Before(timings[j].Span.Start)
	})

	// no timings found process an empty timings correction
	if !(len(timings) > 0) {
		return []Correction{{Span: coll.Span}}
	}

	var res []Correction

	// check prior to the first timing
	if v := timings[0]; v.Start.After(coll.Span.Start) {
		res = append(res, Correction{
			Span: Span{
				Start: coll.Span.Start,
				End:   v.Start,
			},
		})
	}

	// first timing
	res = append(res, Correction{
		Span:   timings[0].Span,
		Timing: &timings[0],
	})

	// subsequent timings, checking for gaps
	for i := 1; i < len(timings); i++ {
		if timings[i].Start.After(timings[i-1].End) {
			res = append(res, Correction{
				Span: Span{
					Start: timings[i-1].End,
					End:   timings[i].Start,
				},
			})
		}
		res = append(res, Correction{
			Span:   timings[i].Span,
			Timing: &timings[i],
		})
	}

	// check after the last timing
	if v := timings[len(timings)-1]; v.End.Before(coll.Span.End) {
		res = append(res, Correction{
			Span: Span{
				Start: v.End,
				End:   coll.Span.End,
			},
		})
	}

	return res
}

// PolarityCorrections returns a slice of Correction values to account for any changes in Polarity.
func (s *Set) PolarityCorrections(coll Collection) []Correction {
	var polarities []Polarity

	// build a slice of associated polarities
	for _, p := range s.Polarities() {
		if p.Station != coll.Stream.Station {
			continue
		}
		if p.Location != coll.Stream.Location {
			continue
		}

		if p.Subsource != coll.Component.Subsource && p.Subsource != "" {
			continue
		}

		if !p.Span.Overlaps(coll.Span) {
			continue
		}

		polarities = append(polarities, p)
	}

	sort.Slice(polarities, func(i, j int) bool {
		return polarities[i].Span.Start.Before(polarities[j].Span.Start)
	})

	// no polarities found process an empty polarity correction
	if !(len(polarities) > 0) {
		return []Correction{{Span: coll.Span}}
	}

	var res []Correction

	// check prior to the first gain
	if v := polarities[0]; v.Start.After(coll.Span.Start) {
		res = append(res, Correction{
			Span: Span{
				Start: coll.Span.Start,
				End:   v.Start,
			},
		})
	}

	res = append(res, Correction{
		Span:     polarities[0].Span,
		Polarity: &polarities[0],
	})

	// subsequent gains, checking for gaps
	for i := 1; i < len(polarities); i++ {
		if polarities[i].Start.After(polarities[i-1].End) {
			res = append(res, Correction{
				Span: Span{
					Start: polarities[i-1].End,
					End:   polarities[i].Start,
				},
			})
		}
		res = append(res, Correction{
			Span:     polarities[i].Span,
			Polarity: &polarities[i],
		})
	}

	// check after the last gain
	if v := polarities[len(polarities)-1]; v.End.Before(coll.Span.End) {
		res = append(res, Correction{
			Span: Span{
				Start: v.End,
				End:   coll.Span.End,
			},
		},
		)
	}

	return res
}

// GainCorrections returns a slice of Correction values to account for any changes in Gain.
func (s *Set) GainCorrections(coll Collection) []Correction {
	var gains []Gain

	for _, g := range s.Gains() {
		if g.Station != coll.Stream.Station {
			continue
		}
		if g.Location != coll.Stream.Location {
			continue
		}
		if g.Subsource != coll.Component.Subsource && g.Subsource != "" {
			continue
		}
		if !g.Span.Overlaps(coll.Span) {
			continue
		}
		gains = append(gains, g)
	}

	sort.Slice(gains, func(i, j int) bool {
		return gains[i].Span.Start.Before(gains[j].Span.Start)
	})

	// no gains found process a nominal correction
	if !(len(gains) > 0) {
		return []Correction{{Span: coll.Span}}
	}

	// gains will not overlap, but there may be gaps
	var res []Correction

	// check prior to the first gain
	if v := gains[0]; v.Start.After(coll.Span.Start) {
		res = append(res, Correction{
			Span: Span{
				Start: coll.Span.Start,
				End:   v.Start,
			},
		})
	}

	// first gain
	res = append(res, Correction{
		Span: gains[0].Span,
		Gain: &gains[0],
	})

	// subsequent gains, checking for gaps
	for i := 1; i < len(gains); i++ {
		if gains[i].Start.After(gains[i-1].End) {
			res = append(res, Correction{
				Span: Span{
					Start: gains[i-1].End,
					End:   gains[i].Start,
				},
			})
		}
		res = append(res, Correction{
			Span: gains[i].Span,
			Gain: &gains[i],
		})
	}

	// check after the last gain
	if v := gains[len(gains)-1]; v.End.Before(coll.Span.End) {
		res = append(res, Correction{
			Span: Span{
				Start: v.End,
				End:   coll.Span.End,
			},
		})
	}

	return res
}

// SensornCalibrationCorrections returns a slice of Correction values to account for any changes in Sensor Calibration.
func (s *Set) SensorCalibrationCorrections(coll Collection) []Correction {

	var sensors []Calibration
	for _, c := range s.Calibrations() {
		if c.Make != coll.InstalledSensor.Make {
			continue
		}
		if c.Model != coll.InstalledSensor.Model {
			continue
		}
		if c.Serial != coll.InstalledSensor.Serial {
			continue
		}
		if !coll.Span.Overlaps(c.Span) {
			continue
		}
		sensors = append(sensors, c)
	}
	sort.Slice(sensors, func(i, j int) bool {
		return sensors[i].Span.Start.Before(sensors[j].Span.Start)
	})

	// no calibrations found return the nominal correction
	if !(len(sensors) > 0) {
		return []Correction{{Span: coll.Span}}
	}

	// calibrations will not overlap, but there may be gaps

	var res []Correction

	// check prior to the first gain
	if v := sensors[0]; v.Start.After(coll.Span.Start) {
		res = append(res, Correction{
			Span: Span{
				Start: coll.Span.Start,
				End:   v.Start,
			},
		})
	}

	// first calibration
	res = append(res, Correction{
		Span:              sensors[0].Span,
		SensorCalibration: &sensors[0],
	})

	// subsequent calibrations, checking for gaps
	for i := 1; i < len(sensors); i++ {
		if sensors[i].Start.After(sensors[i-1].End) {
			res = append(res, Correction{
				Span: Span{
					Start: sensors[i-1].End,
					End:   sensors[i].Start,
				},
			})
		}
		res = append(res, Correction{
			Span:              sensors[i].Span,
			SensorCalibration: &sensors[i],
		})
	}

	// check after the last gain
	if v := sensors[len(sensors)-1]; v.End.Before(coll.Span.End) {
		res = append(res, Correction{
			Span: Span{
				Start: v.End,
				End:   coll.Span.End,
			},
		})
	}

	return res
}

// DataloggerCalibrationCorrections returns a slice of Correction values to account for any changes in DeployedDatalogger Calibration.
func (s *Set) DataloggerCalibrationCorrections(coll Collection) []Correction {

	var dataloggers []Calibration
	for _, c := range s.Calibrations() {
		if c.Make != coll.DeployedDatalogger.Make {
			continue
		}
		if c.Model != coll.DeployedDatalogger.Model {
			continue
		}
		if c.Serial != coll.DeployedDatalogger.Serial {
			continue
		}
		if !coll.Span.Overlaps(c.Span) {
			continue
		}
		dataloggers = append(dataloggers, c)
	}
	sort.Slice(dataloggers, func(i, j int) bool {
		return dataloggers[i].Span.Start.Before(dataloggers[j].Span.Start)
	})

	// no calibrations found return the nominal correction
	if !(len(dataloggers) > 0) {
		return []Correction{{Span: coll.Span}}
	}

	// calibrations will not overlap, but there may be gaps

	var res []Correction
	// check prior to the first gain
	if v := dataloggers[0]; v.Start.After(coll.Span.Start) {
		res = append(res, Correction{
			Span: Span{
				Start: coll.Span.Start,
				End:   v.Start,
			},
		})
	}

	// first calibration
	res = append(res, Correction{
		Span:                  dataloggers[0].Span,
		DataloggerCalibration: &dataloggers[0],
	})

	// subsequent calibrations, checking for gaps
	for i := 1; i < len(dataloggers); i++ {
		if dataloggers[i].Start.After(dataloggers[i-1].End) {
			res = append(res, Correction{
				Span: Span{
					Start: dataloggers[i-1].End,
					End:   dataloggers[i].Start,
				},
			})
		}
		res = append(res, Correction{
			Span:                  dataloggers[i].Span,
			DataloggerCalibration: &dataloggers[i],
		})
	}

	// check after the last gain
	if v := dataloggers[len(dataloggers)-1]; v.End.Before(coll.Span.End) {
		res = append(res, Correction{
			Span: Span{
				Start: v.End,
				End:   coll.Span.End,
			},
		})
	}

	return res
}

// PreampCorrections returns a slice of Correction values to account for any changes in Preamp settings.
func (s *Set) PreampCorrections(coll Collection) []Correction {

	var preamps []Preamp
	for _, p := range s.Preamps() {
		if p.Station != coll.Stream.Station {
			continue
		}
		if p.Location != coll.Stream.Location {
			continue
		}
		if p.Subsource != coll.Component.Subsource && p.Subsource != "" {
			continue
		}
		if !coll.Span.Overlaps(p.Span) {
			continue
		}
		preamps = append(preamps, p)
	}
	sort.Slice(preamps, func(i, j int) bool {
		return preamps[i].Span.Start.Before(preamps[j].Span.Start)
	})

	// no preamps found return an empty span
	if !(len(preamps) > 0) {
		return []Correction{{Span: coll.Span}}
	}

	var res []Correction

	// check prior to the first preamp
	if v := preamps[0]; v.Start.After(coll.Span.Start) {
		res = append(res, Correction{
			Span: Span{
				Start: coll.Span.Start,
				End:   v.Start,
			},
		})
	}

	// first preamp
	res = append(res, Correction{
		Span:   preamps[0].Span,
		Preamp: &preamps[0],
	})

	// subsequent preamps, checking for gaps
	for i := 1; i < len(preamps); i++ {
		if preamps[i].Start.After(preamps[i-1].End) {
			res = append(res, Correction{
				Span: Span{
					Start: preamps[i-1].End,
					End:   preamps[i].Start,
				},
			})
		}
		res = append(res, Correction{
			Span:   preamps[i].Span,
			Preamp: &preamps[i],
		})
	}

	// check after the last gain
	if v := preamps[len(preamps)-1]; v.End.Before(coll.Span.End) {
		res = append(res, Correction{
			Span: Span{
				Start: v.End,
				End:   coll.Span.End,
			},
		})
	}

	return res
}

// TelemetryCorrections returns a slice of Correction values to account for any changes in Preamp settings.
func (s *Set) TelemetryCorrections(coll Collection) []Correction {

	var telemetries []Telemetry
	for _, t := range s.Telemetries() {
		if t.Station != coll.Stream.Station {
			continue
		}
		if t.Location != coll.Stream.Location {
			continue
		}
		if !coll.Span.Overlaps(t.Span) {
			continue
		}
		telemetries = append(telemetries, t)
	}
	sort.Slice(telemetries, func(i, j int) bool {
		return telemetries[i].Span.Start.Before(telemetries[j].Span.Start)
	})

	// no telemetries found return an empty span
	if !(len(telemetries) > 0) {
		return []Correction{{Span: coll.Span}}
	}

	var res []Correction

	// check prior to the first preamp
	if v := telemetries[0]; v.Start.After(coll.Span.Start) {
		res = append(res, Correction{
			Span: Span{
				Start: coll.Span.Start,
				End:   v.Start,
			},
		})
	}

	// first telemetry
	res = append(res, Correction{
		Span:      telemetries[0].Span,
		Telemetry: &telemetries[0],
	})

	// subsequent telemetries, checking for gaps
	for i := 1; i < len(telemetries); i++ {
		if telemetries[i].Start.After(telemetries[i-1].End) {
			res = append(res, Correction{
				Span: Span{
					Start: telemetries[i-1].End,
					End:   telemetries[i].Start,
				},
			})
		}
		res = append(res, Correction{
			Span:      telemetries[i].Span,
			Telemetry: &telemetries[i],
		})
	}

	// check after the last telemetry
	if v := telemetries[len(telemetries)-1]; v.End.Before(coll.Span.End) {
		res = append(res, Correction{
			Span: Span{
				Start: v.End,
				End:   coll.Span.End,
			},
		})
	}

	return res
}
