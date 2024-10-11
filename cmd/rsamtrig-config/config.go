package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var configRe = regexp.MustCompile(`^([A-Z]+)(/[A-Z0-9]+)*(:[0-9]+)*$`)

type Config struct {
	Station  string
	Location string
	Level    float64
}

func NewConfig(str string) (Config, error) {

	parts := configRe.FindStringSubmatch(str)
	if len(parts) != 4 {
		return Config{}, fmt.Errorf("invalid config")
	}

	var location string
	if strings.HasPrefix(parts[2], "/") {
		location = parts[2][1:]
	}

	var level float64
	if strings.HasPrefix(parts[3], ":") {
		l, err := strconv.ParseFloat(parts[3][1:], 64)
		if err != nil {
			return Config{}, err
		}
		level = l
	}

	config := Config{
		Station:  parts[1],
		Location: location,
		Level:    level,
	}

	return config, nil
}
