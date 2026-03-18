package main

import (
	"encoding/xml"
	"io"
	"os"

	stationxml "github.com/GeoNet/delta/internal/stationxml/v1.2"
)

func DecodeStationXML(rd io.Reader) (*stationxml.FDSNStationXML, error) {
	var root stationxml.FDSNStationXML
	if err := xml.NewDecoder(rd).Decode(&root); err != nil {
		return nil, err
	}

	return &root, nil

}

func ReadStationXML(path string) (*stationxml.FDSNStationXML, error) {
	file, err := os.Open(path) //nolint:gosec // disable G304
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	res, err := DecodeStationXML(file)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return res, nil
}

func StationXML(path string) (*stationxml.FDSNStationXML, error) {
	switch {
	case path != "":
		return ReadStationXML(path)
	default:
		return DecodeStationXML(os.Stdin)
	}
}
