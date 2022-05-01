package meta

import (
	"strconv"
	"strings"
)

// ParseInt converts a string into a int value. After trimming spaces, an empty string will return zero,
// otherwise the converted number is returned or an error if the conversion failed.
func ParseInt(s string) (int, error) {
	if s = strings.TrimSpace(s); s != "" {
		return strconv.Atoi(s)
	}
	return 0, nil
}

// ParseFloat64 converts a string into a float64 value. After trimming spaces, an empty string will return zero,
// otherwise the converted number is returned or an error if the conversion failed.
func ParseFloat64(s string) (float64, error) {
	if s = strings.TrimSpace(s); s != "" {
		return strconv.ParseFloat(s, 64)
	}
	return 0.0, nil
}

// ParseSamplingRate converts a string into a float64 value representing a sampling rate. After trimming spaces, an empty string will return zero,
// otherwise the converted sample rate is returned or an error if the float64 conversion failed.
// For initial sampling rates that are found to be negative the inverse of the negative sampling rate is returned,
// the negative number indicates that the value is actually a sampling period.
func ParseSamplingRate(s string) (float64, error) {
	v, err := ParseFloat64(s)
	if err != nil {
		return 0.0, err
	}
	if v < 0.0 {
		return -1.0 / v, nil
	}
	return v, nil
}
