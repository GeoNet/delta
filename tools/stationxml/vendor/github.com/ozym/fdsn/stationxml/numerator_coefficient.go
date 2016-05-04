package stationxml

type NumeratorCoefficient struct {
	Coefficient int32   `xml:"i,attr"`
	Value       float64 `xml:",chardata"`
}
