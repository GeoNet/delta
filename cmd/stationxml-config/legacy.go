package main

import (
	"strings"
)

func LegacyPrivate(code string) bool {
	switch code {
	case "AVIS":
		return true
	default:
		return false
	}
}

func LegacyStorageFormat(datalogger string) string {
	switch {
	case strings.HasPrefix(datalogger, "Centaur"):
		return "Steim1"
	case strings.HasPrefix(datalogger, "Titan"):
		return "Steim1"
	default:
		return "Steim2"
	}
}

func LegacyDescription(datalogger string) string {
	switch {
	case strings.HasPrefix(datalogger, "Q330"):
		return "Q330"
	default:
		return datalogger
	}
}

func LegacyFrequency(code string) float64 {
	switch {
	case strings.HasPrefix(code, "LH"):
		return 0.1
	case strings.HasPrefix(code, "LK"):
		return 0.1
	case strings.HasPrefix(code, "BH"):
		return 1.0
	case strings.HasPrefix(code, "BN"):
		return 1.0
	case strings.HasPrefix(code, "HN"):
		return 1.0
	case strings.HasPrefix(code, "HH"):
		return 1.0
	case strings.HasPrefix(code, "VH"):
		return 0.05
	case strings.HasPrefix(code, "LT"):
		return 0.1
	case strings.HasPrefix(code, "BT"):
		return 1.0
	case strings.HasPrefix(code, "VT"):
		return 0.05
	case strings.HasPrefix(code, "VD"):
		return 0.05
	case strings.HasPrefix(code, "LD"):
		return 0.1
	default:
		return 15.0
	}
}
