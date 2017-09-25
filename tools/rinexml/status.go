package main

import (
	"encoding/xml"
)

type Mark struct {
	XMLName xml.Name `xml:"mark"`

	Code    string `xml:"code,attr"`
	Name    string `xml:"name,attr"`
	Lat     string `xml:"lat,attr"`
	Lon     string `xml:"lon,attr"`
	Opened  string `xml:"opened,attr"`
	Closed  string `xml:"closed,attr,omitempty"`
	Network string `xml:"network,attr"`
}

type Marks struct {
	XMLName xml.Name `xml:"marks"`

	Name  string `xml:"name,attr"`
	Url   string `xml:"url,attr,omitempty"`
	Marks []Mark
}

func (m Marks) Marshal() ([]byte, error) {
	s, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(s, []byte{'\n', '\n'}...)...), nil
}
