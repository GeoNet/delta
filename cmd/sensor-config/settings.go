package main

import (
	"github.com/GeoNet/delta/meta"
)

func (s Settings) Groups(set *meta.Set) Network {

	var groups []Group

	if group, ok := s.InstalledSensors(set, "Air pressure sensor", &s.acoustic, s.networks); ok {
		groups = append(groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Broadband seismometer", &s.seismic, s.networks, "Broadband Seismometer"); ok {
		groups = append(groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Coastal sea level gauge", &s.water, s.coastal); ok {
		groups = append(groups, group)
	}

	if group, ok := s.InstalledSensors(set, "DART bottom pressure recorder", &s.water, s.dart); ok {
		groups = append(groups, group)
	}

	if group, ok := s.EnviroSensor(set, "Environmental sensor", s.enviro); ok {
		groups = append(groups, group)
	}

	if group, ok := s.Doases(set, "DOAS spectrometer", s.doas); ok {
		groups = append(groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Geomagnetic sensor", &s.geomag, s.magnetic); ok {
		groups = append(groups, group)
	}

	if group, ok := s.Gnss(set, "GNSS/GPS", s.gnss); ok {
		groups = append(groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Lake level gauge", &s.water, s.lentic); ok {
		groups = append(groups, group)
	}

	if group, ok := s.ManualCollection(set, "Manual collection", s.manual); ok {
		groups = append(groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Short period seismometer", &s.seismic, s.networks, "Short Period Seismometer"); ok {
		groups = append(groups, group)
	}

	if group, ok := s.InstalledSensors(set, "Strong motion sensor", &s.strong, s.networks); ok {
		groups = append(groups, group)
	}

	if group, ok := s.Cameras(set, "Camera", s.camera); ok {
		groups = append(groups, group)
	}

	return Network{Groups: groups}
}
