package main

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func isAllUpperCase(label string) bool {
	for _, v := range label {
		if !unicode.IsNumber(v) && unicode.ToUpper(v) == v {
			continue
		}
		return false
	}
	return true
}

func cleanLabel(label string) string {
	for _, s := range []string{"/", ".", "+", " "} {
		label = strings.ReplaceAll(label, s, "-")
	}
	return strings.TrimRight(label, "-")
}

func (g Generate) SensorName(model SensorModel, sensor Sensor, label string) string {
	return fmt.Sprintf("sensor_%s_%s", model.Make(), cleanLabel(label))
}

// DataloggerBits is a notional dynamic range for a datalogger assuming a peak to peak range of 20 Volts. Zero is returned if no range can be found.
func (g Generate) DataloggerBits(datalogger Datalogger) int {
	for _, f := range datalogger.Filters {
		m, ok := g.FilterMap[f]
		if !ok {
			continue
		}
		for _, s := range m {
			switch s.Type {
			case "a2d":
				return 1 + int(math.Round(math.Log10(s.Gain*20)/math.Log10(2.0)))
			}
		}
	}
	return 0
}

// DataloggerPreamp returns whether a datalogger which employs FIR filters has an extra PAZ stage.
func (g Generate) DataloggerPreamp(datalogger Datalogger) bool {
	var paz, fir int
	for _, f := range datalogger.Filters {
		m, ok := g.FilterMap[f]
		if !ok {
			continue
		}
		for _, s := range m {
			switch s.Type {
			case "paz":
				paz++
			case "fir":
				fir++
			}
		}
	}
	return paz > 0 && fir > 0
}

func (g Generate) DataloggerName(model DataloggerModel, datalogger Datalogger, label string, rate float64) string {
	if label = strings.Split(cleanLabel(label), "-")[0]; isAllUpperCase(label) {
		label = cases.Title(language.English).String(strings.ToLower(label))
	}

	bits := g.DataloggerBits(datalogger)

	if rate < 0.0 {
		rate = -1.0 / rate
	}

	switch {
	case rate > 1.0:
		if g.DataloggerPreamp(datalogger) {
			return fmt.Sprintf("datalogger_%s_%s-Pre_%dbits_%gsps", model.Make(), label, bits, rate)
		}
		return fmt.Sprintf("datalogger_%s_%s_%dbits_%gsps", model.Make(), label, bits, rate)
	default:
		if g.DataloggerPreamp(datalogger) {
			return fmt.Sprintf("datalogger_%s_%s-Pre_%dbits_%gs", model.Make(), label, bits, 1.0/rate)
		}
		return fmt.Sprintf("datalogger_%s_%s_%dbits_%gs", model.Make(), label, bits, 1.0/rate)
	}
}
