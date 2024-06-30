package main

import (
	"fmt"
	"strings"

	"github.com/GeoNet/delta/meta"
)

func (s Settings) Network(set *meta.Set) (Network, error) {

	var network Network

	// this is the legacy mechanism - which has trouble with sensors at sites which have unrelated networks
	for _, n := range strings.Split(s.networks, ",") {
		if err := network.InstalledSensors(set, &s.combined, strings.TrimSpace(n), ""); err != nil {
			return Network{}, fmt.Errorf("unable to build seismic details (%s): %v", n, err)
		}
	}

	if err := network.InstalledSensors(set, &s.water, s.coastal, "Coastal"); err != nil {
		return Network{}, fmt.Errorf("unable to build water details: %v", err)
	}

	if err := network.InstalledSensors(set, &s.water, s.lentic, "Lake"); err != nil {
		return Network{}, fmt.Errorf("unable to build water details: %v", err)
	}

	if err := network.EnviroSensor(set, s.enviro, "Environmental Sensor"); err != nil {
		return Network{}, fmt.Errorf("unable to build envirosensor configuration: %v", err)
	}

	if err := network.Dart(set, s.dart, "DART Bottom Pressure Recorder"); err != nil {
		return Network{}, fmt.Errorf("unable to build dart configuration: %v", err)
	}

	if err := network.Gnss(set, "GNSS Antenna", "GNSS Receiver"); err != nil {
		return Network{}, fmt.Errorf("unable to build gnss configuration: %v", err)
	}

	if err := network.Camera(set, s.volcano, "Camera"); err != nil {
		return Network{}, fmt.Errorf("unable to build camera configuration: %v", err)
	}

	if err := network.Camera(set, s.building, "Camera"); err != nil {
		return Network{}, fmt.Errorf("unable to build camera configuration: %v", err)
	}

	if err := network.Doas(set, s.doas, "DOAS"); err != nil {
		return Network{}, fmt.Errorf("unable to build camera configuration: %v", err)
	}

	if err := network.ManualCollection(set, s.manual, "Manual Collection"); err != nil {
		return Network{}, fmt.Errorf("unable to build camera configuration: %v", err)
	}

	return network, nil
}

func (s Settings) Groups(set *meta.Set) Network {

	var network Network

	if group, ok := s.InstalledSensors(set, "Air pressure sensor", &s.acoustic, s.networks); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Broadband seismometer", &s.seismic, s.networks, "Broadband Seismometer"); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Coastal sea level gauge", &s.water, s.coastal); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.InstalledSensors(set, "DART bottom pressure recorder", &s.water, s.dart); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.EnviroSensor(set, "Environmental sensor", s.enviro); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.Doases(set, "DOAS spectrometer", s.doas); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Geomagnetic sensor", &s.geomag, s.magnetic); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.Gnss(set, "GNSS/GPS", s.gnss); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Lake level gauge", &s.water, s.lentic); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.ManualCollection(set, "Manual collection", s.manual); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Short period seismometer", &s.seismic, s.networks, "Short Period Seismometer"); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Strong motion sensor", &s.strong, s.networks); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.Cameras(set, "Camera", s.camera); ok {
		network.Groups = append(network.Groups, group)
	}

	return network
}
