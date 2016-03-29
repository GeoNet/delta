package stationxml

import "fmt"

// This type contains a URI and description for external data that users may want to reference in StationXML.
type ExternalReference struct {
	URI         string
	Description string
}

func (e ExternalReference) IsValid() error {

	if !(len(e.URI) > 0) {
		return fmt.Errorf("empty external reference uri")
	}
	if !(len(e.Description) > 0) {
		return fmt.Errorf("empty external reference description")
	}

	return nil
}
