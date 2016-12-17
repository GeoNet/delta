package metadb

import (
	"time"

	"github.com/GeoNet/delta/meta"
)

type Installation struct {
	Station    string
	Location   string
	Sensor     meta.InstalledSensor
	Datalogger meta.DeployedDatalogger
	Start      time.Time
	End        time.Time
}

func (m *MetaDB) Installations(station string) ([]Installation, error) {
	var installations []Installation

	recorders, err := m.StationInstalledRecorders(station)
	if err != nil {
		return nil, err
	}
	for _, recorder := range recorders {
		installations = append(installations, Installation{
			Station:  station,
			Location: recorder.Location,
			Sensor:   recorder.InstalledSensor,
			Datalogger: meta.DeployedDatalogger{
				Install: meta.Install{
					Equipment: meta.Equipment{
						Model: recorder.DataloggerModel,
					},
					Span: meta.Span{
						Start: recorder.Start,
						End:   recorder.End,
					},
				},
			},
			Start: recorder.Start,
			End:   recorder.End,
		})
	}

	connections, err := m.StationConnections(station)
	if err != nil {
		return nil, err
	}

	for _, connection := range connections {
		location, err := m.Site(station, connection.Location)
		if err != nil {
			return nil, err
		}
		if location == nil {
			return nil, nil
		}

		installs, err := m.StationInstalledSensors(station)
		if err != nil {
			return nil, err
		}

		for _, sensorInstall := range installs {

			deploys, err := m.ConnectionInstalledSensorDeployedDataloggers(connection, sensorInstall)
			if err != nil {
				return nil, err
			}
			if deploys == nil {
				continue
			}
			for _, dataloggerDeploy := range deploys {

				start, end := connection.Start, connection.End
				if sensorInstall.Start.After(start) {
					start = sensorInstall.Start
				}
				if dataloggerDeploy.Start.After(start) {
					start = dataloggerDeploy.Start
				}
				if sensorInstall.End.Before(end) {
					end = sensorInstall.End
				}
				if dataloggerDeploy.End.Before(end) {
					end = dataloggerDeploy.End
				}

				installations = append(installations, Installation{
					Station:    station,
					Location:   connection.Location,
					Sensor:     sensorInstall,
					Datalogger: dataloggerDeploy,
					Start:      start,
					End:        end,
				})
			}
		}
	}

	return installations, nil
}
