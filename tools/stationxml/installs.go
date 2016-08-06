package main

import (
	"path/filepath"
	"sort"

	"github.com/GeoNet/delta/meta"
)

func SensorInstalls(install string) (map[string]meta.InstalledSensorList, error) {
	sensorInstalls := make(map[string]meta.InstalledSensorList)

	// build sensor installation details
	var sensors meta.InstalledSensorList
	if err := meta.LoadList(filepath.Join(install, "sensors.csv"), &sensors); err != nil {
		return nil, err
	}
	for _, s := range sensors {
		if _, ok := sensorInstalls[s.StationCode]; ok {
			sensorInstalls[s.StationCode] = append(sensorInstalls[s.StationCode], s)
		} else {
			sensorInstalls[s.StationCode] = meta.InstalledSensorList{s}
		}
	}

	var recorders meta.InstalledRecorderList
	if err := meta.LoadList(filepath.Join(install, "recorders.csv"), &recorders); err != nil {
		return nil, err
	}
	for _, r := range recorders {
		if _, ok := sensorInstalls[r.StationCode]; ok {
			sensorInstalls[r.StationCode] = append(sensorInstalls[r.StationCode], r.InstalledSensor)
		} else {
			sensorInstalls[r.StationCode] = meta.InstalledSensorList{r.InstalledSensor}
		}
	}

	var gauges meta.InstalledGaugeList
	if err := meta.LoadList(filepath.Join(install, "gauges.csv"), &gauges); err != nil {
		return nil, err
	}
	for _, g := range gauges {
		s := meta.InstalledSensor{
			Install: meta.Install{
				Equipment: meta.Equipment{
					Make:   g.Make,
					Model:  g.Model,
					Serial: g.Serial,
				},
				Span: meta.Span{
					Start: g.Start,
					End:   g.End,
				},
			},
			Offset: meta.Offset{
				Height: g.Height,
				North:  g.North,
				East:   g.East,
			},
			Orientation: meta.Orientation{
				Dip:     g.Dip,
				Azimuth: g.Azimuth,
			},
			StationCode:  g.StationCode,
			LocationCode: g.LocationCode,
		}

		if _, ok := sensorInstalls[s.StationCode]; ok {
			sensorInstalls[s.StationCode] = append(sensorInstalls[s.StationCode], s)
		} else {
			sensorInstalls[s.StationCode] = meta.InstalledSensorList{s}
		}
	}

	// sort each sensor install
	for i, _ := range sensorInstalls {
		sort.Sort(sensorInstalls[i])
	}

	return sensorInstalls, nil
}
