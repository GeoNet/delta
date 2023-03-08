package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Frequency represents the Response frequency to use for a given Channel code prefix.
type Frequency struct {
	Prefix string
	Value  float64
}

// NewFrequency builds a Frequency based on an encoded text string of the form <prefix>:<value>
func NewFrequency(str string) (Frequency, error) {
	parts := strings.SplitN(str, ":", 2)
	switch n := len(parts); {
	case n > 1:
		v, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return Frequency{}, err
		}
		return Frequency{Prefix: strings.TrimSpace(parts[0]), Value: v}, nil
	case n > 0:
		v, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return Frequency{}, err
		}
		return Frequency{Value: v}, nil
	default:
		return Frequency{}, nil
	}
}

func (f Frequency) String() string {
	if len(f.Prefix) > 0 {
		return fmt.Sprintf("%s:%g", f.Prefix, f.Value)
	}
	return fmt.Sprintf("%g", f.Value)
}

func (f Frequency) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *Frequency) UnmarshalText(data []byte) error {
	v, err := NewFrequency(string(data))
	if err != nil {
		return err
	}
	*f = v
	return nil
}

type Frequencies map[string]float64

func NewFrequencies(freqs ...Frequency) Frequencies {
	res := make(map[string]float64)
	for _, f := range freqs {
		res[f.Prefix] = f.Value
	}
	return res
}

// Get gets the value associated with the given key.
// If there are no values associated with the key, Get
// returns zero and false, otherwise it returns the value and true.
func (f Frequencies) Get(key string) (float64, bool) {
	if f == nil {
		return 0.0, false
	}
	v, ok := f[key]
	return v, ok
}

// Set sets the key to value. It replaces any existing values.
func (f Frequencies) Set(key string, value float64) {
	f[key] = value
}

// Del deletes the value associated with key.
func (f Frequencies) Del(key string) {
	delete(f, key)
}

// Has checks whether a given key is set.
func (f Frequencies) Has(key string) bool {
	_, ok := f[key]
	return ok
}

// Find checks the code against the prefixes and returns a value with the longest prefixes. For multiple matches the returned value will be matched alphabetically.
// If no matches are found then zero and false are returned, otherwise, the value and true will be returned.
func (f Frequencies) Find(code string) (float64, bool) {
	match := -1

	var key string
	var value float64
	for k, v := range f {
		if !strings.HasPrefix(code, k) {
			continue
		}
		if len(k) < match {
			continue
		}
		if len(k) == match && k < key {
			continue
		}

		match, key = len(k), k
		value = v
	}

	if match < 0 {
		return 0.0, false
	}

	return value, true
}
