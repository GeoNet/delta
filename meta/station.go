package meta

import (
	"time"
)

type Station struct {
	Network   string    `csv:"Network Code",`
	Code      string    `csv:"Station Code",`
	Name      string    `csv:"Station Name",`
	Latitude  float64   `csv:"Latitude",`
	Longitude float64   `csv:"Longitude",`
	Depth     float64   `csv:"Depth",`
	StartTime time.Time `csv:"Start Time"`
	EndTime   time.Time `csv:"End Time"`
}

type Stations []Station

func (s Stations) list() {}
