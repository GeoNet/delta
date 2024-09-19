package main

import (
	"fmt"
	"strings"
)

const stationTemplate = `###
### Delivered by puppet
###
global:{{.Global.Key}}
rsamtrig:{{.Style}}
`

func ToKey(network, station string) string {
	return fmt.Sprintf("%s_%s", network, station)
}

type Station struct {
	Global   Global
	RsamTrig RsamTrig
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

func (s Station) Style() string {
	return fmt.Sprintf("%s_%s", s.RsamTrig.Style(), strings.ToLower(s.Code))
}

func (s Station) Path() string {
	return fmt.Sprintf("station_%s", s.Key())
}
