package stationxml

type Coefficient struct {
	Number uint32  `xml:"number,attr"`
	Value  float64 `xml:",chardata"`
}
