package main

import (
	"fmt"
)

const DefaultFilter = "RMHP(10)>>BW(4,2,15)>>STALTA(0.5,20)"
const DefaultCorrection = -0.05

const autopickTemplate = `###
###
### Delivered by puppet
###
# Defines the filter to be used for picking.
detecFilter = "{{.Filter}}"

# The time correction applied to a detected pick.
timeCorr = {{.Correction}}

# Defines whether or not the streams are picked or not
detecEnable = {{.Enable}}
`

type AutoPick struct {
	Style      string
	Filter     string
	Correction float64
}

func (s AutoPick) Enable() bool {
	switch s.Style {
	case "broadband", "weak":
		return true
	default:
		return false
	}
}

func (s AutoPick) Id() string {
	return "scautopick"
}
func (s AutoPick) Key() string {
	return s.Style
}
func (s AutoPick) Template() string {
	return autopickTemplate
}

func (s AutoPick) Path() string {
	return fmt.Sprintf("profile_%s", s.Key())
}
