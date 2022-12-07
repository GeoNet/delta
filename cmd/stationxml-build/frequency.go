package main

import (
	"strings"
)

func DefaultFrequency(code string) float64 {
	switch {
	case strings.HasPrefix(code, "VH"):
		return 0.05
	case strings.HasPrefix(code, "VT"):
		return 0.05
	case strings.HasPrefix(code, "VD"):
		return 0.05
	case strings.HasPrefix(code, "LH"):
		return 0.1
	case strings.HasPrefix(code, "LT"):
		return 0.1
	case strings.HasPrefix(code, "LK"):
		return 0.1
	case strings.HasPrefix(code, "LD"):
		return 0.1
	case strings.HasPrefix(code, "BH"):
		return 1.0
	case strings.HasPrefix(code, "BN"):
		return 1.0
	case strings.HasPrefix(code, "BT"):
		return 1.0
	case strings.HasPrefix(code, "HN"):
		return 1.0
	case strings.HasPrefix(code, "HH"):
		return 1.0
	default:
		return 15.0
	}
}
