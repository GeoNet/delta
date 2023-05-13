package main

import (
	"fmt"
	"path/filepath"
)

const DefaultFilter = "RMHP(10)>>BW(4,2,15)>>STALTA(0.5,20)"
const DefaultCorrection = -0.05

const autopickTemplate = `###
### Warning, this file is automaticaly generated and may be overwritten.
###
# Defines the filter to be used for picking.
detecFilter = "{{.Filter}}"

# The time correction applied to a detected pick.
timeCorr = {{.Correction}}

# Defines whether or not the streams are picked or not
detecEnable = {{.Enable}}
`

type AutoPick struct {
	Style        string
	Filter       string
	Correction   float64
	Location     string
	SamplingRate float64
}

func (a AutoPick) Enable() bool {
	switch a.Style {
	case "broadband", "weak":
		return true
	default:
		return false
	}
}

func (a AutoPick) Id() string {
	return "scautopick"
}
func (a AutoPick) Key() string {
	return a.Style
}
func (a AutoPick) Template() string {
	return autopickTemplate
}

func (a AutoPick) Path() string {
	return filepath.Join(a.Id(), fmt.Sprintf("profile_%s", a.Key()))
}

func (a AutoPick) Store(path string) error {
	return Store(a, path)
}
