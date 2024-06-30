package main

import (
	"fmt"
)

const stationTemplate = `###
### Delivered by puppet
###
global:{{.Global.Key}}
scautopick:{{.AutoPick.Style}}
`

func ToKey(network, station string) string {
	return fmt.Sprintf("%s_%s", network, station)
}

func ToStream(network, station, location string) string {
	return fmt.Sprintf("%s_%s_%s", network, station, location)
}

type Station struct {
	Global   Global
	AutoPick AutoPick
	Network  string
	Code     string
}

func (s Station) Id() string {
	return "station"
}

func (s Station) Template() string {
	return stationTemplate
}

func (s Station) Key() string {
	return ToKey(s.Network, s.Code)
}

func (s Station) Path() string {
	return fmt.Sprintf("station_%s", s.Key())
}
