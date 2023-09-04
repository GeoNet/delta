package main

import (
	"bytes"
	"strings"
	"time"
)

type Deployment struct {
	Buoy       string
	Deployment string
	Name       string
	Region     string

	Pid      string   // WMO Id A WMO ID (or station id) E.g. 55048.
	Banks    []string // Payload Ids A group of 4 uppercase alphanumeric characters that identifies the payloads from a Tsunami platform.
	Serial   string   // Paroscientific SN Pressure sensor SN in BPR E.g. 129002.
	Platform string   // Platform Type A type, e.g. DART II or DART 4G that determine the formats of messages for data processing E.g. DART 4G

	Latitude  float64 // BPR Drop Position latitude where the BPR is dropped at sea.  e.g. latitude in decimal degree.
	Longitude float64 // BPR Drop Position longitude where the BPR is dropped at sea.  e.g. longitude in decimal degree.
	Depth     float64 // Water Depth Ship surveyed water depth at the BPR drop location, e.g. 4300 meters.

	Start time.Time // Deployment Start Estimated or planned deployment start date and time
	End   time.Time
}

func (d Deployment) Payload(buoy string) string {
	var list []string
	for _, s := range d.Banks {
		list = append(list, buoy+s)
	}
	return strings.Join(list, " ")
}

type Deployments []Deployment

func (d Deployments) Process() ([]byte, error) {

	var res bytes.Buffer
	if err := tmpl.Execute(&res, d); err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}
