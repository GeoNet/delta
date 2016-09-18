package main

import (
	"path/filepath"
	"sort"
	"time"

	"github.com/GeoNet/delta/meta"
)

type Meta struct {
	Networks    map[string]meta.Network
	Stations    map[string]meta.Station
	Connections map[string]meta.ConnectionList
	Sites       map[string]map[string]meta.Site
	Streams     map[string]map[string]meta.StreamList
	Installs    map[string]meta.InstalledSensorList
	Deploys     map[string]meta.DeployedDataloggerList
}

func NewMeta(network, install string) (*Meta, error) {
	networkMap := make(map[string]meta.Network)

	var n meta.NetworkList
	if err := meta.LoadList(filepath.Join(network, "networks.csv"), &n); err != nil {
		return nil, err
	}

	for _, v := range n {
		networkMap[v.Code] = v
	}

	stationMap := make(map[string]meta.Station)

	var stations meta.StationList
	if err := meta.LoadList(filepath.Join(network, "stations.csv"), &stations); err != nil {
		return nil, err
	}

	for _, v := range stations {
		stationMap[v.Code] = v
	}

	connectionMap := make(map[string]meta.ConnectionList)

	var connections meta.ConnectionList
	if err := meta.LoadList(filepath.Join(install, "connections.csv"), &connections); err != nil {
		return nil, err
	}

	for _, c := range connections {
		if _, ok := connectionMap[c.Station]; ok {
			connectionMap[c.Station] = append(connectionMap[c.Station], c)
		} else {
			connectionMap[c.Station] = meta.ConnectionList{c}
		}
	}

	var recorders meta.InstalledRecorderList
	if err := meta.LoadList(filepath.Join(install, "recorders.csv"), &recorders); err != nil {
		return nil, err
	}
	for _, r := range recorders {
		c := meta.Connection{
			Station:  r.Station,
			Location: r.Location,
			Span: meta.Span{
				Start: r.Start,
				End:   r.End,
			},
			Place: r.Station,
			Role:  r.Location,
		}
		if _, ok := connectionMap[c.Station]; ok {
			connectionMap[c.Station] = append(connectionMap[c.Station], c)
		} else {
			connectionMap[c.Station] = meta.ConnectionList{c}
		}
	}

	siteMap := make(map[string]map[string]meta.Site)

	var locations meta.SiteList
	if err := meta.LoadList(filepath.Join(network, "sites.csv"), &locations); err != nil {
		return nil, err
	}

	for _, l := range locations {
		if _, ok := siteMap[l.Station]; !ok {
			siteMap[l.Station] = make(map[string]meta.Site)
		}
		siteMap[l.Station][l.Location] = l
	}

	streamMap := make(map[string]map[string]meta.StreamList)

	var streams meta.StreamList
	if err := meta.LoadList(filepath.Join(install, "streams.csv"), &streams); err != nil {
		return nil, err
	}

	for _, s := range streams {
		if _, ok := streamMap[s.Station]; !ok {
			streamMap[s.Station] = make(map[string]meta.StreamList)
		}

		if _, ok := streamMap[s.Station][s.Location]; ok {
			streamMap[s.Station][s.Location] = append(streamMap[s.Station][s.Location], s)
		} else {
			streamMap[s.Station][s.Location] = meta.StreamList{s}
		}
	}

	sensorInstalls := make(map[string]meta.InstalledSensorList)

	// build sensor installation details
	var sensors meta.InstalledSensorList
	if err := meta.LoadList(filepath.Join(install, "sensors.csv"), &sensors); err != nil {
		return nil, err
	}
	for _, s := range sensors {
		if _, ok := sensorInstalls[s.Station]; ok {
			sensorInstalls[s.Station] = append(sensorInstalls[s.Station], s)
		} else {
			sensorInstalls[s.Station] = meta.InstalledSensorList{s}
		}
	}

	for _, r := range recorders {
		if _, ok := sensorInstalls[r.Station]; ok {
			sensorInstalls[r.Station] = append(sensorInstalls[r.Station], r.InstalledSensor)
		} else {
			sensorInstalls[r.Station] = meta.InstalledSensorList{r.InstalledSensor}
		}
	}

	for i, _ := range sensorInstalls {
		sort.Sort(sensorInstalls[i])
	}

	dataloggerDeploys := make(map[string]meta.DeployedDataloggerList)

	// where the dataloggers were deployed
	var loggers meta.DeployedDataloggerList
	if err := meta.LoadList(filepath.Join(install, "dataloggers.csv"), &loggers); err != nil {
		return nil, err
	}
	for _, d := range loggers {
		if _, ok := dataloggerDeploys[d.Place]; ok {
			dataloggerDeploys[d.Place] = append(dataloggerDeploys[d.Place], d)
		} else {
			dataloggerDeploys[d.Place] = meta.DeployedDataloggerList{d}
		}
	}

	for _, r := range recorders {
		d := meta.DeployedDatalogger{
			Install: meta.Install{
				Equipment: meta.Equipment{
					Make:   r.Make,
					Model:  r.DataloggerModel,
					Serial: r.Serial,
				},
				Span: meta.Span{
					Start: r.Start,
					End:   r.End,
				},
			},
			Place: r.Station,
			Role:  r.Location,
		}
		if _, ok := dataloggerDeploys[d.Place]; ok {
			dataloggerDeploys[d.Place] = append(dataloggerDeploys[d.Place], d)
		} else {
			dataloggerDeploys[d.Place] = meta.DeployedDataloggerList{d}
		}
	}

	// sort each datalogger deployment
	for i, _ := range dataloggerDeploys {
		sort.Sort(dataloggerDeploys[i])
	}

	return &Meta{
		Networks:    networkMap,
		Stations:    stationMap,
		Connections: connectionMap,
		Sites:       siteMap,
		Streams:     streamMap,
		Installs:    sensorInstalls,
		Deploys:     dataloggerDeploys,
	}, nil
}

func (m Meta) GetNetwork(network string) *meta.Network {
	if v, ok := m.Networks[network]; ok {
		return &v
	}
	return nil
}

func (m Meta) GetStation(station string) *meta.Station {
	if v, ok := m.Stations[station]; ok {
		return &v
	}
	return nil
}

func (m Meta) GetStationKeys() []string {
	var keys []string
	for k, _ := range m.Stations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (m Meta) GetStreams(station, location string) meta.StreamList {
	s, ok := m.Streams[station]
	if !ok {
		return nil
	}
	l, ok := s[location]
	if !ok {
		return nil
	}
	return l
}

func (m Meta) GetConnections(station string) meta.ConnectionList {
	if v, ok := m.Connections[station]; ok {
		return v
	}
	return nil
}

func (m Meta) GetSites(station string) map[string]meta.Site {
	if v, ok := m.Sites[station]; ok {
		return v
	}
	return nil
}

func (m Meta) GetSite(station, location string) *meta.Site {
	s := m.GetSites(station)
	if s == nil {
		return nil
	}
	if v, ok := s[location]; ok {
		return &v
	}
	return nil
}

func (m Meta) GetStream(station, location string, sampleRate float64, at time.Time) *meta.Stream {

	var stream *meta.Stream

	for _, s := range m.GetStreams(station, location) {
		if s.SamplingRate != sampleRate {
			continue
		}
		if s.End.Before(at) {
			continue
		}
		if s.Start.After(at) {
			break
		}
		stream = &s
	}

	return stream
}

func (m Meta) GetInstalls(station string) meta.InstalledSensorList {
	if v, ok := m.Installs[station]; ok {
		return v
	}
	return nil
}

func (m Meta) GetDeploys(place string) meta.DeployedDataloggerList {
	if v, ok := m.Deploys[place]; ok {
		return v
	}
	return nil
}
