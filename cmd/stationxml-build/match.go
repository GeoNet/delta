package main

import (
	"github.com/GeoNet/delta/meta"
)

type Match struct {
	externals map[string][]string
	networks  map[string]meta.Network
	stations  map[string]meta.Station
	sites     map[string][]meta.Site
}

func NewMatch(set *meta.Set, external, network, station, location Matcher) Match {
	m := Match{
		externals: make(map[string][]string),
		networks:  make(map[string]meta.Network),
		stations:  make(map[string]meta.Station),
		sites:     make(map[string][]meta.Site),
	}

	for _, n := range set.Networks() {
		if !external.MatchString(n.External) {
			continue
		}
		m.externals[n.External] = append(m.externals[n.External], n.Code)
	}

	for _, n := range set.Networks() {
		if !network.MatchString(n.Code) {
			continue
		}
		if _, ok := m.externals[n.External]; !ok {
			continue
		}

		m.networks[n.Code] = n
	}

	for _, s := range set.Stations() {
		if !station.MatchString(s.Code) {
			continue
		}
		if _, ok := m.networks[s.Network]; !ok {
			continue
		}

		m.stations[s.Code] = s
	}
	for _, s := range set.Sites() {
		if !location.MatchString(s.Location) {
			continue
		}
		if _, ok := m.stations[s.Station]; !ok {
			continue
		}

		m.sites[s.Station] = append(m.sites[s.Station], s)
	}

	return m
}

func (m Match) Externals() map[string][]string {
	externals := make(map[string][]string)
	for k, v := range m.externals {
		externals[k] = append([]string{}, v...)
	}
	return externals
}

func (m Match) Stations() map[string]meta.Station {
	stations := make(map[string]meta.Station)
	for k, v := range m.stations {
		stations[k] = v
	}
	return stations
}

func (m Match) Network(code string) (meta.Network, bool) {
	n, ok := m.networks[code]
	return n, ok
}

func (m Match) Sites(code string) []meta.Site {
	return append([]meta.Site{}, m.sites[code]...)
}
