package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

const DefaultName = "f1-4"
const DefaultFilter = "BW(4,1,4)"

const rsamTrigTemplate = `###
### Delivered by puppet
###
# Defines a list of filters to apply.
rsamFilters = {{.Name}}

# Waveform filter string to apply before RSAM calculation.
rsamFilter.{{.Name}}.filter = {{.Filter}}

# Minimum amplitude level. Used as starting point and lower barrier.
rsamFilter.{{.Name}}.baseLevel = {{.Base}}
`

type RsamTrig struct {
	Station      string
	Name         string
	Filter       string
	Base         float64
	Location     string
	SamplingRate float64
}

func (r RsamTrig) Id() string {
	return "rsamtrig"
}
func (r RsamTrig) Style() string {
	return "rsam"
}
func (r RsamTrig) Key() string {
	return fmt.Sprintf("%s_%s", r.Style(), strings.ToLower(r.Station))
}
func (r RsamTrig) Template() string {
	return rsamTrigTemplate
}

func (r RsamTrig) Path() string {
	return filepath.Join(r.Id(), fmt.Sprintf("profile_%s", r.Key()))
}
