package main

import (
	"fmt"
	"strings"
)

const stationTemplate = `###
### Warning, this file is automaticaly generated and may be overwritten.
###
global:{{.Global.Key}}
scautopick:{{.AutoPick.Style}}
`

type Station struct {
	Global   Global
	AutoPick AutoPick
	Network  string
	Code     string
}

func NewStation(s string) Station {
	r := strings.Split(strings.ToUpper(s), "_")
	switch n := len(r); {
	case n > 1:
		return Station{
			Network: r[0],
			Code:    r[1],
		}
	case n == 1:
		return Station{
			Network: "*",
			Code:    r[0],
		}
	default:
		return Station{
			Network: "*",
			Code:    "*",
		}
	}
}

func (s Station) Id() string {
	return "station"
}

func (s Station) Template() string {
	return stationTemplate
}

func (s Station) Key() string {
	return fmt.Sprintf("%s_%s", s.Network, s.Code)
}

func (s Station) Path() string {
	return fmt.Sprintf("station_%s", s.Key())
}

func (s Station) Store(path string) error {
	return Store(s, path)
}
