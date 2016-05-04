package stationxml

// Network allows grouping of station metadata.
//
// This type represents the Network layer, all station metadata is contained within this element.
// The official name of the network or other descriptive information can be included in the
// Description element. The Network can contain 0 or more Stations.
type Network struct {
	BaseNode

	// The total number of stations contained in this network, including inactive or terminated stations.
	TotalNumberStations uint32 `xml:",omitempty" json:",omitempty"`

	// The total number of stations in this network that were selected by the query that produced this document,
	// even if the stations do not appear in the document. (This might happen if the user only wants a document
	// that goes contains only information at the Network level.)
	SelectedNumberStations uint32 `xml:",omitempty" json:",omitempty"`

	Stations []Station `xml:"Station,omitempty" json:",omitempty"`
}

func (n Network) IsValid() error {

	if err := Validate(n.BaseNode); err != nil {
		return err
	}

	for _, s := range n.Stations {
		if err := Validate(s); err != nil {
			return err
		}
	}

	return nil
}
