package main

import (
	"strings"
)

func Frequency(code string) float64 {
	switch {
	case strings.HasPrefix(code, "V"):
		return 0.05
	case strings.HasPrefix(code, "L"):
		return 0.1
	case strings.HasPrefix(code, "B"):
		return 1.0
	case strings.HasPrefix(code, "H"):
		return 1.0
	case strings.HasPrefix(code, "S"):
		return 15.0
	case strings.HasPrefix(code, "E"):
		return 15.0
	default:
		return 15.0
	}
}
