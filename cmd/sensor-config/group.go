package main

import (
	"github.com/GeoNet/delta/meta"
)

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

	if group, ok := s.InstalledSensors(set, "DART Bottom Pressure Recorder", &s.water, s.dart); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.Doases(set, "DOAS Spectrometer", s.doas); ok {
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

	if group, ok := s.ManualCollection(set, "Manual Collection", s.manual); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Short period seismometer", &s.seismic, s.networks, "Short Period Seismometer"); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Strong motion sensor", &s.strong, s.networks); ok {
		network.Groups = append(network.Groups, group)
	}

	if group, ok := s.Cameras(set, "Volcano Cameras", s.volcano); ok {
		network.Groups = append(network.Groups, group)
	}

	return network
}
